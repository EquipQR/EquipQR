package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

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

func ComparePasswordHash(encodedHash, password string, cfg Argon2Config) bool {
	parts := split(encodedHash, "$")
	if len(parts) != 2 {
		return false
	}

	salt, err1 := base64.RawStdEncoding.DecodeString(parts[0])
	_, err2 := base64.RawStdEncoding.DecodeString(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}

	newHash := argon2.IDKey([]byte(password), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLen)
	return base64.RawStdEncoding.EncodeToString(newHash) == parts[1]
}

func split(input, sep string) []string {
	var parts []string
	i := 0
	for {
		j := i + len(sep)
		if j > len(input) {
			break
		}
		if input[i:j] == sep {
			parts = append(parts, input[:i], input[j:])
			break
		}
		i++
	}
	if len(parts) == 0 {
		return []string{input}
	}
	return parts
}
