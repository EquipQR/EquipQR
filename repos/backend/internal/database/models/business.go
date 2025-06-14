package models

type Business struct {
	ID           string `gorm:"primaryKey"`
	BusinessName string `gorm:"uniqueIndex;size:64;not null"`
}
