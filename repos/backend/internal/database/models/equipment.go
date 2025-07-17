package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Equipment struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BusinessID uuid.UUID      `gorm:"type:uuid;not null;index" json:"businessId"`
	Status     string         `gorm:"type:text;not null;check:status IN ('in service','not in service')" json:"status"`
	Type       string         `gorm:"type:text;not null" json:"type"`
	Location   string         `gorm:"type:text" json:"location"`
	MoreFields datatypes.JSON `gorm:"type:jsonb" json:"moreFields"`
	Business   Business       `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE" json:"business"`
}
