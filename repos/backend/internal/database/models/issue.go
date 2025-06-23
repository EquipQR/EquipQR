package models

import (
	"time"

	"github.com/google/uuid"
)

type Issue struct {
	ID            uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title         string     `gorm:"type:varchar(128);not null" json:"title"`
	EquipmentID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"equipment_id"`
	Description   string     `gorm:"type:text;not null" json:"description"`
	Progress      string     `gorm:"type:text;not null" json:"progress"`
	AssigneeID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"assignee_id"`
	DateSubmitted time.Time  `gorm:"not null" json:"date_submitted"`
	DateCompleted *time.Time `gorm:"default:null" json:"date_completed,omitempty"`
	Equipment     Equipment  `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment"`
	Assignee      User       `gorm:"foreignKey:AssigneeID;constraint:OnDelete:RESTRICT" json:"assignee"`
}
