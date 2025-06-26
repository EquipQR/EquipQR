package utils

import (
	"fmt"
	"strings"
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
	if tokenString == "" {
		return "", fmt.Errorf("empty JWT string")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		fmt.Println("❌ JWT parse error:", err)
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token or claims")
	}

	return claims.UserID, nil
}

func GetJWTSecret() []byte {
	secret := strings.TrimSpace(AppConfig.JWT_Secret)
	fmt.Printf("🔐 GetJWTSecret() returning (sanitized): %q\n", secret)
	return []byte(secret)
}

func SetOrRemoveSessionCookie(c *fiber.Ctx, token string) {
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
		cookie.Expires = time.Now().Add(time.Duration(AppConfig.Cookie_Expiry_Days) * 24 * time.Hour)
	}

	c.Cookie(cookie)
}

func ValidateJWTFromCookie(c *fiber.Ctx) (string, error) {

	fmt.Println("JWT secret in app:", string(GetJWTSecret()))

	cookie := c.Cookies("session")
	fmt.Println("🍪 Session cookie read:", cookie)
	parts := strings.Split(cookie, ".")
	if len(parts) != 3 {
		fmt.Println("⚠️ Token malformed:", cookie)
	}
	if cookie == "" {
		fmt.Println("⚠️ No session cookie found")
		return "", fiber.ErrUnauthorized
	}

	userID, err := ValidateJWT(cookie)
	if err != nil {
		fmt.Println("❌ JWT validation failed:", err)
		return "", fiber.ErrUnauthorized
	}

	fmt.Println("✅ JWT valid, userID:", userID)
	return userID, nil
}
