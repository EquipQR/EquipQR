package models

import (
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/gorm"
)

type UserBusiness struct {
	ID         string   `gorm:"primaryKey;type:varchar(64);not null"`
	UserID     string   `gorm:"type:uuid;not null;index"`
	BusinessID string   `gorm:"type:uuid;not null;index"`
	User       User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE"`
}

func (user_business *UserBusiness) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	user_business.ID = utils.PikaGenerator.NextID("user_business")
	return nil
}
