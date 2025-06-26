package models

import (
	"github.com/google/uuid"
)

type Credential struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"`
	CredentialID    []byte    `gorm:"uniqueIndex;not null"`
	PublicKey       []byte    `gorm:"not null"`
	AttestationType string    `gorm:"not null"`
	AAGUID          []byte    `gorm:"not null"`
	SignCount       uint32    `gorm:"not null"`
	CloneWarning    bool      `gorm:"not null"`
}
