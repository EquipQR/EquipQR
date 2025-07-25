package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/EquipQR/equipqr/backend/internal/middleware"
	"github.com/EquipQR/equipqr/backend/internal/s3"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func RegisterMediaRoutes(app *fiber.App) {
	// Single file upload (generic, no issue ID)
	app.Post("/api/upload", middleware.RequireUser, UploadFile)

	// Multi-file upload for a specific issue
	app.Post("/api/issue/:id/attachments", middleware.RequireUser, UploadFiles)

	// File download by key
	app.Get("/files/:key", middleware.RequireUser, GetFile)
}

func UploadFile(c *fiber.Ctx) error {
	config := utils.LoadConfigFromEnv()

	bucket := config.MinioBucket
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "file not provided")
	}

	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)

	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "could not open file")
	}
	defer file.Close()

	_, err = s3.Client.PutObject(context.Background(), bucket, objectName, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "upload failed")
	}

	return c.JSON(fiber.Map{
		"key": objectName,
		"url": fmt.Sprintf("/files/%s", objectName),
	})
}

func UploadFiles(c *fiber.Ctx) error {
	config := utils.LoadConfigFromEnv()

	bucket := config.MinioBucket

	// ⛳ Extract issue ID from route
	issueIDStr := c.Params("id")
	issueID, err := uuid.Parse(issueIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid issue ID")
	}
	log.Printf("📎 Attaching files to issue: %s", issueID)

	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid multipart form")
	}

	files := form.File["files"]
	if len(files) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "no files provided")
	}

	type UploadedFile struct {
		ID       uuid.UUID `json:"id"`
		Key      string    `json:"key"`
		URL      string    `json:"url"`
		FileName string    `json:"file_name"`
		IssueID  uuid.UUID `json:"issue_id"`
	}

	var uploaded []UploadedFile

	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("could not open file: %s", fileHeader.Filename))
		}
		defer src.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, src); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("could not read file: %s", fileHeader.Filename))
		}
		log.Printf("🧪 Uploading file: %s (%d bytes)", fileHeader.Filename, buf.Len())

		objectKey := fmt.Sprintf("%d-%s", time.Now().UnixNano(), fileHeader.Filename)
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		_, err = s3.Client.PutObject(context.Background(), bucket, objectKey, bytes.NewReader(buf.Bytes()), int64(buf.Len()), minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			log.Printf("❌ MinIO upload error: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("upload failed for %s: %v", fileHeader.Filename, err))
		}

		uploaded = append(uploaded, UploadedFile{
			ID:       uuid.New(),
			Key:      objectKey,
			URL:      fmt.Sprintf("/files/%s", objectKey),
			FileName: fileHeader.Filename,
			IssueID:  issueID,
		})
	}

	log.Printf("✅ Uploaded %d file(s) for issue %s", len(uploaded), issueID)

	return c.JSON(fiber.Map{
		"uploaded": uploaded,
	})
}

func GetFile(c *fiber.Ctx) error {
	config := utils.LoadConfigFromEnv()
	bucket := config.MinioBucket
	key := c.Params("key")

	obj, err := s3.Client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "file not found")
	}
	defer obj.Close()

	stat, err := obj.Stat()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "file stat error")
	}

	c.Type(stat.ContentType)
	c.Response().Header.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", key))
	c.Response().Header.Set("Content-Length", fmt.Sprintf("%d", stat.Size))

	_, err = io.Copy(c.Response().BodyWriter(), obj)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "stream error")
	}

	return nil
}
