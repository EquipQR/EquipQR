package handlers

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/EquipQR/equipqr/backend/internal/s3"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
)

var config = utils.LoadConfigFromEnv()

func UploadFile(c *fiber.Ctx) error {
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

func GetFile(c *fiber.Ctx) error {
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
