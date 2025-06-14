package models

type UserBusiness struct {
	ID         string   `gorm:"primaryKey"`
	UserID     string   `gorm:"type:uuid;not null;index"`
	BusinessID string   `gorm:"type:uuid;not null;index"`
	User       User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Business   Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE"`
}
