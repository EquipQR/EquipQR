package models

import (
	"time"

	"github.com/google/uuid"
)

type PendingJoinRequest struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null" json:"business_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
