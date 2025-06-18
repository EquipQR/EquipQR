package models

import (
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Equipment struct {
	ID         string         `gorm:"primaryKey;type:varchar(64);not null" json:"id"`
	BusinessID string         `gorm:"type:uuid;not null;index" json:"businessId"`
	Status     string         `gorm:"type:text;not null;check:status IN ('in service','not in service')" json:"status"`
	Type       string         `gorm:"type:text;not null" json:"type"`
	Location   string         `gorm:"type:text" json:"location"`
	MoreFields datatypes.JSON `gorm:"type:jsonb" json:"moreFields"`
	Business   Business       `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
}

func (equipment *Equipment) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	equipment.ID = utils.PikaGenerator.NextID("equipment")
	return nil
}
