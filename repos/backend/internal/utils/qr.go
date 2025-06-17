package utils

import (
	"github.com/skip2/go-qrcode"
)

func GenerateQRCodeBytes(equipmentID string) ([]byte, error) {
	return qrcode.Encode(equipmentID, qrcode.Medium, 256)
}
