package repositories

import (
	"archive/zip"
	"bytes"
	"fmt"

	"github.com/EquipQR/equipqr/backend/internal/utils"
)

func GenerateQRCodeZipBytes(equipmentIDs []string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for _, id := range equipmentIDs {
		equipment, err := GetEquipmentByID(id)
		if err != nil {
			return nil, fmt.Errorf("equipment with ID %s not found: %w", id, err)
		}

		qrCode, err := utils.GenerateQRCodeBytes(id)
		if err != nil {
			return nil, fmt.Errorf("failed to generate QR code for %s: %w", id, err)
		}

		fileWriter, err := zipWriter.Create(fmt.Sprintf("%s.png", equipment.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to create zip entry for %s: %w", id, err)
		}

		_, err = fileWriter.Write(qrCode)
		if err != nil {
			return nil, fmt.Errorf("failed to write QR for %s: %w", id, err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip: %w", err)
	}

	return buf.Bytes(), nil
}

func GenerateSingleQRCodeBytes(equipmentID string) ([]byte, string, error) {
	equipment, err := GetEquipmentByID(equipmentID)
	if err != nil {
		return nil, "", fmt.Errorf("equipment with ID %s not found: %w", equipmentID, err)
	}

	qrCode, err := utils.GenerateQRCodeBytes(equipmentID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate QR for %s: %w", equipmentID, err)
	}

	return qrCode, equipment.ID.String(), nil
}
