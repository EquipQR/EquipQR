package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Email    string    `gorm:"size:100;not null" json:"email"`
	Password string    `gorm:"size:256;not null" json:"-"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
}
