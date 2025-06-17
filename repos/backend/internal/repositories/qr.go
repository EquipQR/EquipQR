package repositories

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"

	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func GenerateQRCodeZip(context *fiber.Ctx) error {
	// Declare a struct to hold the input JSON
	type RequestBody struct {
		EquipmentIDs []string `json:"equipment_ids"`
	}

	// Parse the request body into a RequestBody struct
	var requestBody RequestBody
	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if the EquipmentIDs array is empty
	if len(requestBody.EquipmentIDs) == 0 {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No equipment IDs provided",
		})
	}

	// Create a buffer to write the zip file to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Iterate over equipment IDs and add QR codes to the zip file
	for _, id := range requestBody.EquipmentIDs {
		// Validate the equipment ID by checking if it exists in the database
		equipment, err := GetEquipmentByID(id)
		if err != nil {
			return context.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("Equipment with ID %s not found: %v", id, err),
			})
		}

		// Generate the QR code as a byte slice
		qrCode, err := utils.GenerateQRCodeBytes(id)
		if err != nil {
			return fmt.Errorf("failed to generate QR code for %s: %v", id, err)
		}

		// Create a new zip file entry for each QR code
		fileWriter, err := zipWriter.Create(fmt.Sprintf("%s.png", equipment.ID))
		if err != nil {
			return fmt.Errorf("failed to create zip file entry for %s: %v", id, err)
		}

		// Write the QR code image to the zip file
		_, err = fileWriter.Write(qrCode)
		if err != nil {
			return fmt.Errorf("failed to write QR code for %s to zip file: %v", id, err)
		}
	}

	// Close the zip writer
	err := zipWriter.Close()
	if err != nil {
		return fmt.Errorf("failed to close zip file: %v", err)
	}

	// Set the response headers and return the zip file
	context.Set("Content-Type", "application/zip")
	context.Set("Content-Disposition", "attachment; filename=qr_codes.zip")

	// Write the buffer (zip file) to the response
	err = context.Send(buf.Bytes())
	if err != nil {
		log.Printf("Error sending file: %v", err)
		return err
	}

	return nil
}

func GenerateSingleQRCode(context *fiber.Ctx) error {
	type RequestBody struct {
		EquipmentID string `json:"equipment_id"`
	}

	var requestBody RequestBody
	if err := context.BodyParser(&requestBody); err != nil {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if requestBody.EquipmentID == "" {
		return context.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No equipment ID provided",
		})
	}

	// Validate equipment id: check if it exists in the database
	equipment, err := GetEquipmentByID(requestBody.EquipmentID)
	if err != nil {
		return context.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": fmt.Sprintf("Equipment with ID %s not found: %v", requestBody.EquipmentID, err),
		})
	}

	qrCode, err := utils.GenerateQRCodeBytes(requestBody.EquipmentID)
	if err != nil {
		return fmt.Errorf("failed to generate QR code for %s: %v", requestBody.EquipmentID, err)
	}

	context.Set("Content-Type", "image/png")
	context.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.png", equipment.ID))

	// Write QRCode image as a response
	err = context.Send(qrCode)
	if err != nil {
		log.Printf("Error sending the file: %v", err)
		return err
	}

	return nil
}
