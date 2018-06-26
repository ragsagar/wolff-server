package model

import (
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/satori/go.uuid"
)

func HashPassword(rawPassword string) string {
	hasher := sha256.New()
	io.WriteString(hasher, rawPassword)
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateUUID() string {
	id := uuid.NewV4()
	return id.String()
}
