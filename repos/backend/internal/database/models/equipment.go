package models

import (
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Equipment struct {
	ID         string         `gorm:"primaryKey;type:varchar(64);not null"`
	BusinessID string         `gorm:"type:uuid;not null;index"`
	Status     string         `gorm:"type:text;not null;check:status IN ('in service','not in service')"`
	Type       string         `gorm:"type:text;not null"`
	Location   string         `gorm:"type:text"`
	MoreFields datatypes.JSON `gorm:"type:jsonb"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE"`
}

func (equipment *Equipment) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	equipment.ID = utils.PikaGenerator.NextID("equipment")
	return nil
}
