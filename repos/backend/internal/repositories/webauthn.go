package repositories

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/EquipQR/equipqr/backend/internal/database"
	"github.com/EquipQR/equipqr/backend/internal/database/models"
	"github.com/EquipQR/equipqr/backend/internal/utils"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

var WebAuthn *webauthn.WebAuthn

func InitWebAuthn() {
	authn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: "EquipQR",
		RPID:          "equipqr.io",
		RPOrigin:      "https://equipqr.io",
	})
	if err != nil {
		panic(err)
	}
	WebAuthn = authn
}

func BeginWebAuthnRegistration(userID string, c *fiber.Ctx) (*protocol.CredentialCreation, error) {
	user, err := GetUserByIDWithCredentials(userID)
	if err != nil {
		return nil, err
	}

	opts, session, err := WebAuthn.BeginRegistration(user)
	if err != nil {
		return nil, err
	}

	saveWebAuthnSession(c, user.ID.String(), session)
	return opts, nil
}

func FinishWebAuthnRegistration(c *fiber.Ctx) error {
	userID := getUserIDFromSession(c)
	user, err := GetUserByIDWithCredentials(userID)
	if err != nil {
		return err
	}

	session := loadWebAuthnSession(c, user.ID.String())
	if session == nil {
		return errors.New("no webauthn session found")
	}

	req := convertFiberRequestToHTTP(c)
	cred, err := WebAuthn.FinishRegistration(user, *session, req)
	if err != nil {
		return err
	}

	return SaveWebAuthnCredential(user.ID, cred)
}

func BeginWebAuthnLogin(email string, c *fiber.Ctx) (*protocol.CredentialAssertion, error) {
	log.Printf("[WebAuthn] Begin login for email: %s", email)

	user, err := GetUserByEmailWithCredentials(email)
	if err != nil {
		log.Printf("[WebAuthn] Failed to find user or load credentials: %v", err)
		return nil, err
	}

	if len(user.Credentials) == 0 {
		log.Printf("[WebAuthn] User %s has no WebAuthn credentials", user.Email)
	}

	opts, session, err := WebAuthn.BeginLogin(user)
	if err != nil {
		log.Printf("[WebAuthn] Failed to generate login options: %v", err)
		return nil, err
	}

	log.Printf("[WebAuthn] Generated login options for user %s", user.Email)

	saveWebAuthnSession(c, user.ID.String(), session)
	log.Printf("[WebAuthn] Session saved for user ID %s", user.ID.String())
	log.Printf("[WebAuthn] Session saved for user ID %s", user.ID.String())

	return opts, nil
}

func FinishWebAuthnLogin(c *fiber.Ctx) (string, error) {
	userID := getUserIDFromSession(c)
	user, err := GetUserByIDWithCredentials(userID)
	if err != nil {
		return "", err
	}

	session := loadWebAuthnSession(c, user.ID.String())
	if session == nil {
		return "", errors.New("no webauthn session found")
	}

	req := convertFiberRequestToHTTP(c)
	_, err = WebAuthn.FinishLogin(user, *session, req)
	if err != nil {
		return "", err
	}

	return utils.GenerateJWT(user.ID.String())
}

func SaveWebAuthnCredential(userID uuid.UUID, cred *webauthn.Credential) error {
	model := models.Credential{
		UserID:          userID,
		CredentialID:    cred.ID,
		PublicKey:       cred.PublicKey,
		AttestationType: cred.AttestationType,
		AAGUID:          cred.Authenticator.AAGUID,
		SignCount:       cred.Authenticator.SignCount,
		CloneWarning:    cred.Authenticator.CloneWarning,
	}
	return database.DB.Create(&model).Error
}

func GetUserByIDWithCredentials(id string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Credentials").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmailWithCredentials(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Preload("Credentials").First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func saveWebAuthnSession(c *fiber.Ctx, key string, session *webauthn.SessionData) {
	data, _ := json.Marshal(session)
	c.Cookie(&fiber.Cookie{
		Name:  "webauthn_session_" + key,
		Value: string(data),
		Path:  "/",
	})
}

func loadWebAuthnSession(c *fiber.Ctx, key string) *webauthn.SessionData {
	val := c.Cookies("webauthn_session_" + key)
	if val == "" {
		return nil
	}
	var session webauthn.SessionData
	_ = json.Unmarshal([]byte(val), &session)
	return &session
}

func getUserIDFromSession(c *fiber.Ctx) string {
	userID, err := utils.ValidateJWTFromCookie(c)
	if err != nil {
		return ""
	}
	return userID
}

func convertFiberRequestToHTTP(c *fiber.Ctx) *http.Request {
	req := new(http.Request)
	_ = fasthttpadaptor.ConvertRequest(c.Context(), req, true)
	return req
}
