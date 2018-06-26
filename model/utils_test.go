package model

import (
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()
	if uuid == "" {
		t.Errorf("Failed to generate uuid. Returned empty string.")
	}
}

func TestHashPassword(t *testing.T) {
	password := "password"
	expected := "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
	got := HashPassword(password)
	if got != expected {
		t.Errorf("%s not equal to received %s", expected, got)
	}
}
