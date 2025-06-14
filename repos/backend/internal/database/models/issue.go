package models

import "time"

type Issue struct {
	ID            string     `gorm:"primaryKey"`
	EquipmentID   string     `gorm:"not null;index"`
	Description   string     `gorm:"type:text;not null"`
	Progress      string     `gorm:"type:text;not null"`
	AssigneeID    string     `gorm:"not null;index"`
	DateSubmitted time.Time  `gorm:"not null"`
	DateCompleted *time.Time `gorm:"default:null"`
	Equipment Equipment `gorm:"foreignKey:EquipmentID;constraint:OnDelete:CASCADE"`
	Assignee  User      `gorm:"foreignKey:AssigneeID;constraint:OnDelete:SET NULL"`
}
