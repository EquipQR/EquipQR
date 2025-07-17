package models

import (
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username    string       `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Email       string       `gorm:"size:100;not null" json:"email"`
	Password    string       `gorm:"size:256;not null" json:"-"`
	IsActive    bool         `gorm:"default:true" json:"is_active"`
	Credentials []Credential `gorm:"foreignKey:UserID"`
}

func (u *User) WebAuthnID() []byte {
	return []byte(u.ID.String())
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Username
}

func (u *User) WebAuthnIcon() string {
	return "" // optional, can return a URL to profile image
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	var creds []webauthn.Credential
	for _, c := range u.Credentials {
		creds = append(creds, webauthn.Credential{
			ID:              c.CredentialID,
			PublicKey:       c.PublicKey,
			AttestationType: c.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:       c.AAGUID,
				SignCount:    c.SignCount,
				CloneWarning: c.CloneWarning,
			},
		})
	}
	return creds
}
