package models

import (
	"time"

	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/gorm"
)

type Issue struct {
	ID          string `gorm:"primaryKey;type:varchar(64);not null" json:"id"`
	Title       string `gorm:"type:varchar(128);not null" json:"title"`
	EquipmentID string `gorm:"not null;index" json:"equipment_id"`
	Description string `gorm:"type:text;not null" json:"description"`
	Progress    string `gorm:"type:text;not null" json:"progress"`
	AssigneeID  string `gorm:"not null;index" json:"assignee_id"`

	DateSubmitted time.Time  `gorm:"not null" json:"date_submitted"`
	DateCompleted *time.Time `gorm:"default:null" json:"date_completed,omitempty"`

	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE" json:"equipment"`
	Assignee  User      `gorm:"foreignKey:AssigneeID;constraint:OnDelete:RESTRICT" json:"assignee"`
}

func (issue *Issue) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	issue.ID = utils.PikaGenerator.NextID("issue")
	issue.DateSubmitted = time.Now().UTC()
	return nil
}
