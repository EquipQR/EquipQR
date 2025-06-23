package models

import "github.com/google/uuid"

type Business struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BusinessName string    `json:"businessName" validate:"required,min=2,max=64"`
}
