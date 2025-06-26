package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var AppConfig = LoadConfigFromEnv()

type Claims struct {
	UserID string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	claims := Claims{
		UserID: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(AppConfig.JWT_Expiry_Minutes) * time.Minute)),
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
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", err
	}

	return claims.UserID, nil
}

func GetJWTSecret() []byte {
	return []byte(AppConfig.JWT_Secret)
}

func SetOrRemoveSessionCookie(c *fiber.Ctx, token string) {
	cookie := &fiber.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	if token == "" {
		cookie.Expires = time.Now().Add(-1 * time.Hour)
	} else {
		cookie.Expires = time.Now().Add(time.Duration(AppConfig.Cookie_Expiry_Days) * 24 * time.Hour)
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
