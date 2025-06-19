package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Config struct {
	Memory  uint32
	Time    uint32
	Threads uint8
	KeyLen  uint32
	SaltLen uint32
}

var DefaultArgon2Config = Argon2Config{
	Memory:  64 * 1024, // 64 MB
	Time:    3,
	Threads: 2,
	KeyLen:  32,
	SaltLen: 16,
}

func GeneratePasswordHash(password string, cfg Argon2Config) (string, error) {
	salt := make([]byte, cfg.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLen)

	encoded := base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)
	return encoded, nil
}

func ComparePasswordHash(encodedHash string, password string, cfg Argon2Config) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false, errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, errors.New("invalid base64 salt")
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, errors.New("invalid base64 hash")
	}

	newHash := argon2.IDKey([]byte(password), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLen)
	if len(newHash) != len(expectedHash) {
		return false, errors.New("hash length mismatch")
	}

	for i := range newHash {
		if newHash[i] != expectedHash[i] {
			return false, nil
		}
	}

	return true, nil
}
