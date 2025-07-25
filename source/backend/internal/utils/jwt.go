package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	config := LoadConfigFromEnv()
	claims := Claims{
		UserID: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.JWT_Expiry_Minutes) * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTSecret())

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (string, error) {
	if tokenString == "" {
		return "", fmt.Errorf("empty JWT string")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token or claims")
	}

	return claims.UserID, nil
}

func GetJWTSecret() []byte {
	config := LoadConfigFromEnv()
	secret := strings.TrimSpace(config.JWT_Secret)
	return []byte(secret)
}

func SetOrRemoveSessionCookie(c *fiber.Ctx, token string) {
	config := LoadConfigFromEnv()
	cookie := &fiber.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	}

	if token == "" {
		cookie.Expires = time.Now().Add(-1 * time.Hour)
	} else {
		cookie.Expires = time.Now().Add(time.Duration(config.Cookie_Expiry_Days) * 24 * time.Hour)
	}

	c.Cookie(cookie)
}

func ValidateJWTFromCookie(c *fiber.Ctx) (string, error) {

	cookie := c.Cookies("session")

	if cookie == "" {
		return "", fiber.ErrUnauthorized
	}

	userID, err := ValidateJWT(cookie)
	if err != nil {
		return "", fiber.ErrUnauthorized
	}

	return userID, nil
}
