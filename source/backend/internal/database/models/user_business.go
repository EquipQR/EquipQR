package models

import "github.com/google/uuid"

type UserBusiness struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	BusinessID uuid.UUID `gorm:"type:uuid;not null;index" json:"business_id"`
	IsAdmin    bool      `gorm:"default:false" json:"is_admin"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`
	Business   Business  `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
}
