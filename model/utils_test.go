package model

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	if uuid == "" {
		t.Errorf("Failed to generate uuid. Returned empty string.")
	}
}

func TestHashPassword(t *testing.T) {
	password := "password"
	gotHash, _ := HashPassword(password)
	err := bcrypt.CompareHashAndPassword([]byte(gotHash), []byte(password))
	if err != nil {
		t.Errorf("Hashes are not matching.")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "password"
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		t.Fatal(err)
	}
	hashedPassword := string(bytes)
	if !CheckPasswordHash(password, hashedPassword) {
		t.Errorf("Checking password hash failed.")
	}
	if CheckPasswordHash("password2", hashedPassword) {
		t.Errorf("CheckPasswordHash returned false positive.")
	}
}

func TestGenerateToken(t *testing.T) {
	token1 := GenerateTokenKey()
	token2 := GenerateTokenKey()
	if token1 == token2 {
		t.Errorf("Two tokens can't be same.")
	}
	if len(token1) != 40 {
		t.Errorf("Length of token should be 40.")
	}
	// t.Errorf("Token2: %s", token1)
}

func TestGenerateRandomData(t *testing.T) {
	r1, err := GenerateRandomData(12)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := GenerateRandomData(20)
	if err != nil {
		t.Fatal(err)
	}
	if len(r1) != 12 || len(r2) != 20 {
		t.Errorf("Length of the random data is not correct.")
	}
}
