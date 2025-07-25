package utils

import (
	"github.com/google/uuid"
)

func MustParseUUID(id string) uuid.UUID {
	parsed, err := uuid.Parse(id)
	if err != nil {
		panic("invalid UUID: " + id)
	}
	return parsed
}
