package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;size:64;not null"`
	Email    string `gorm:"size:100;not null"`
	Password string `gorm:"size:256;not null"`
	IsActive bool   `gorm:"default:true"`
}
