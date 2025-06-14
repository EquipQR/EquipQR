package models

import (
	"gorm.io/datatypes"
)

type Equipment struct {
	ID          string         `gorm:"primaryKey"`
	BusinessID  string         `gorm:"type:uuid;not null;index"`
	Status      string         `gorm:"type:text;not null;check:status IN ('in service','not in service')"`
	Type        string         `gorm:"type:text;not null"`
	EquipmentID string         `gorm:"type:uuid;not null;uniqueIndex"`
	Location    string         `gorm:"type:text"`
	MoreFields  datatypes.JSON `gorm:"type:jsonb"`

	Business Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE"`
}
