package repositories

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
		RPID:          "app.equipqr.io", // equipqr.io
		RPOrigin:      "https://app.equipqr.io:7878",
	})
	if err != nil {
		panic(err)
	}
	WebAuthn = authn
}

func BeginWebAuthnRegistration(userID string, c *fiber.Ctx) (*protocol.CredentialCreation, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is empty")
	}

	user, err := GetUserByIDWithCredentials(userID)
	if err != nil {
		fmt.Printf("failed to get user for WebAuthn registration: %v\n", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	opts, session, err := WebAuthn.BeginRegistration(user)
	if err != nil {
		fmt.Printf("WebAuthn.BeginRegistration failed: %v\n", err)
		return nil, fmt.Errorf("begin registration failed: %w", err)
	}

	saveWebAuthnSession(c, user.ID.String(), session)
	return opts, nil
}

func FinishWebAuthnRegistration(c *fiber.Ctx) error {
	userID := getUserIDFromSession(c)
	if userID == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	user, err := GetUserByIDWithCredentials(userID)
	if err != nil {
		fmt.Printf("❌ failed to fetch user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch user",
		})
	}

	session := loadWebAuthnSession(c, user.ID.String())
	if session == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no webauthn session found",
		})
	}

	req := convertFiberRequestToHTTP(c)
	cred, err := WebAuthn.FinishRegistration(user, *session, req)
	if err != nil {
		fmt.Printf("❌ WebAuthn finish error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "registration failed",
		})
	}

	if err := SaveWebAuthnCredential(user.ID, cred); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save credential",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func BeginWebAuthnLogin(email string, c *fiber.Ctx) (*protocol.CredentialAssertion, error) {
	log.Printf("[WebAuthn] Begin login for email: %s", email)

	user, err := GetUserByEmailWithCredentials(email)
	if err != nil {
		log.Printf("[WebAuthn] Failed to find user or load credentials: %v", err)
		return nil, err
	}
	for _, cred := range user.Credentials {
		fmt.Printf("[WebAuthn] → Allowed credential for %s: %x\n", user.Email, cred.ID)
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

	return opts, nil
}

func FinishWebAuthnLogin(c *fiber.Ctx) (string, error) {
	type Assertion struct {
		ID                     string         `json:"id"`
		RawID                  string         `json:"rawId"`
		Type                   string         `json:"type"`
		Response               map[string]any `json:"response"`
		Email                  string         `json:"email"`
		ClientExtensionResults map[string]any `json:"clientExtensionResults"`
	}

	var input Assertion
	if err := c.BodyParser(&input); err != nil || input.Email == "" || input.ID == "" {
		fmt.Println("[WebAuthn] ❌ Login failed: invalid or missing data in request body")
		return "", errors.New("invalid request")
	}

	fmt.Printf("[WebAuthn] 🔐 Login attempt for email: %s\n", input.Email)

	user, err := GetUserByEmail(input.Email)
	if err != nil {
		fmt.Printf("[WebAuthn] ❌ Login failed: user not found: %v\n", err)
		return "", fmt.Errorf("user lookup failed: %w", err)
	}

	session := loadWebAuthnSession(c, user.ID.String())
	if session == nil {
		fmt.Printf("[WebAuthn] ❌ Login failed: no session found for user %s\n", user.ID.String())
		return "", errors.New("no webauthn session found")
	}

	req := convertFiberRequestToHTTP(c)
	decodedID, _ := base64.RawURLEncoding.DecodeString(input.ID)

	_, err = WebAuthn.FinishLogin(user, *session, req)
	if err != nil {
		fmt.Printf("[WebAuthn] ❌ FinishLogin failed: %v\n", err)

		return "", fmt.Errorf(
			"webauthn login failed for user %s: %w (submitted credential: %x, allowed: %d credential(s))",
			user.ID.String(), err, decodedID, len(session.AllowedCredentialIDs),
		)
	}

	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		fmt.Printf("[WebAuthn] ❌ Failed to generate JWT: %v\n", err)
		return "", fmt.Errorf("token generation failed: %w", err)
	}

	utils.SetOrRemoveSessionCookie(c, token)

	fmt.Printf("[WebAuthn] ✅ Login successful for user %s\n", user.ID.String())
	return token, nil
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
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	var user models.User
	err = database.DB.Preload("Credentials").First(&user, "id = ?", parsedID).Error
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
