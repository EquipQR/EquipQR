package models

import (
	"time"

	"github.com/google/uuid"
)

type PendingJoinRequest struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
