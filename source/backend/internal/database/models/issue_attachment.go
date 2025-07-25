package models

import (
	"time"

	"github.com/google/uuid"
)

type IssueAttachment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	IssueID   uuid.UUID `gorm:"type:uuid;not null;index" json:"issue_id"`
	FileKey   string    `gorm:"type:text;not null" json:"file_key"`          // S3 key/path
	FileName  string    `gorm:"type:varchar(255);not null" json:"file_name"` // Original filename
	MimeType  string    `gorm:"type:varchar(128);not null" json:"mime_type"`
	Uploaded  time.Time `gorm:"autoCreateTime" json:"uploaded"`

	Issue Issue `gorm:"foreignKey:IssueID;constraint:OnDelete:CASCADE" json:"-"`
}