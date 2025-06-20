package models

import (
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey;type:varchar(64);not null" json:"id"`
	Username string `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Email    string `gorm:"size:100;not null" json:"email"`
	Password string `gorm:"size:256;not null" json:"-"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	if utils.PikaGenerator == nil {
		panic("PikaGenerator is not initialized")
	}
	user.ID = utils.PikaGenerator.NextID("user")
	return nil
}
