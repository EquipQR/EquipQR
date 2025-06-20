package utils

import (
	"github.com/skip2/go-qrcode"
)

// GenerateQRCodeBytes generates a PNG-encoded QR code as a byte slice from the given equipment ID.
//
// The QR code is generated with medium error correction level and a size of 256x256 pixels.
// Returns the PNG bytes or an error if encoding fails.
func GenerateQRCodeBytes(equipmentID string) ([]byte, error) {
	return qrcode.Encode(equipmentID, qrcode.Medium, 256)
}
