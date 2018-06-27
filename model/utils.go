package model

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/satori/go.uuid"
)

// HashPassword: Returns sha256 hash of the given password.
func HashPassword(rawPassword string) string {
	hasher := sha256.New()
	io.WriteString(hasher, rawPassword)
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateUUID: Generate and return a uuid4 string.
func GenerateUUID() string {
	id := uuid.NewV4()
	return id.String()
}
