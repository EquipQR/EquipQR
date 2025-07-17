package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Business struct {
	ID              uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	BusinessName    string         `json:"businessName" validate:"required,min=2,max=64"`
	BusinessEmail   string         `gorm:"size:100" json:"businessEmail" validate:"required,email"`
	Phone           string         `gorm:"size:32" json:"phone"`
	CountryCode     string         `gorm:"size:8" json:"countryCode"`
	Type            string         `json:"type" validate:"required,oneof=aviation transport mechanics factory manufacturing construction mining 'oil & gas' marine automotive healthcare hospitality other"`
	CompanySize     string         `gorm:"size:64" json:"companySize"`
	Country         string         `gorm:"size:64" json:"country"`
	UserCanRegister bool           `json:"userCanRegister" gorm:"default:true"`
	LoginMethods    pq.StringArray `gorm:"type:text[]" json:"loginMethods" validate:"dive,oneof=password magic_link oauth"`
}
