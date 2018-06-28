package model

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns sha256 hash of the given password.
func HashPassword(password string) (string, error) {
	// hasher := sha256.New()
	// io.WriteString(hasher, rawPassword)
	// return hex.EncodeToString(hasher.Sum(nil))
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash returns true if the hash given password matches with given hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateUUID returns a uuid4 string.
func GenerateUUID() string {
	id := uuid.NewV4()
	return id.String()
}

// GenerateTokenKey returns the key of length 40 that can be used in Token model.
func GenerateTokenKey() string {
	randomData, _ := GenerateRandomData(20)
	return hex.EncodeToString(randomData)
}

// GenerateRandomData returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomData(n int) ([]byte, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return bytes, nil
}
