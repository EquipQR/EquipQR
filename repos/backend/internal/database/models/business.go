package models

import (
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/gorm"
)

type Business struct {
	ID           string `gorm:"primaryKey;type:varchar(64);not null"`
	BusinessName string `gorm:"size:64;not null"`
}

func (business *Business) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	business.ID = utils.PikaGenerator.NextID("business")
	return nil
}
