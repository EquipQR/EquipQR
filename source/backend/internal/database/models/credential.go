package models

import (
	"encoding/hex"
	"encoding/json"

	"github.com/duo-labs/webauthn/webauthn"
	"github.com/google/uuid"
)

type Credential struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CredentialID    []byte    `gorm:"uniqueIndex;not null" json:"-"`
	PublicKey       []byte    `gorm:"not null" json:"-"`
	AttestationType string    `gorm:"not null" json:"attestation_type"`
	AAGUID          []byte    `gorm:"not null" json:"-"`
	SignCount       uint32    `gorm:"not null" json:"sign_count"`
	CloneWarning    bool      `gorm:"not null" json:"clone_warning"`
}

func (c Credential) ToWebAuthn() webauthn.Credential {
	return webauthn.Credential{
		ID:              c.CredentialID,
		PublicKey:       c.PublicKey,
		AttestationType: c.AttestationType,
		Authenticator: webauthn.Authenticator{
			AAGUID:       c.AAGUID,
			SignCount:    c.SignCount,
			CloneWarning: c.CloneWarning,
		},
	}
}

func (c Credential) MarshalJSON() ([]byte, error) {
	type Alias Credential
	return json.Marshal(&struct {
		CredentialID string `json:"credential_id"`
		PublicKey    string `json:"public_key"`
		AAGUID       string `json:"aaguid"`
		*Alias
	}{
		CredentialID: hex.EncodeToString(c.CredentialID),
		PublicKey:    hex.EncodeToString(c.PublicKey),
		AAGUID:       hex.EncodeToString(c.AAGUID),
		Alias:        (*Alias)(&c),
	})
}
