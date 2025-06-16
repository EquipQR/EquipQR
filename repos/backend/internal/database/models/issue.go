package models

import (
	"time"

	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/gorm"
)

type Issue struct {
	ID            string     `gorm:"primaryKey;type:varchar(64);not null"`
	EquipmentID   string     `gorm:"not null;index"`
	Description   string     `gorm:"type:text;not null"`
	Progress      string     `gorm:"type:text;not null"`
	AssigneeID    string     `gorm:"not null;index"`
	DateSubmitted time.Time  `gorm:"not null"`
	DateCompleted *time.Time `gorm:"default:null"`
	Equipment     Equipment  `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE"`
	Assignee      User       `gorm:"foreignKey:AssigneeID;constraint:OnDelete:SET NULL"`
}

func (issue *Issue) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	issue.ID = utils.PikaGenerator.NextID("issue")
	return nil
}
