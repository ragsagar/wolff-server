package model

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUser(t *testing.T) {
	password := "password123"
	user := User{
		ID:    "5d6e34c8-46b7-11e6-ba7c-cafec0ffee12",
		Email: "testuser3@gmail.com",
	}
	err := user.SetPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPasswordHash(password, user.Password) {
		t.Errorf("Password hashes are not matching.")
	}

	userStr := fmt.Sprintf("%s", user)
	if userStr != "User<testuser3@gmail.com>" {
		t.Errorf("String function didn't return correct value %s", userStr)
	}

	expected, _ := json.Marshal(user)
	userJSON, err := user.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(userJSON) != string(expected) {
		fmt.Println(string(userJSON))
		t.Errorf("Return json didn't match with expected json %s", string(userJSON))
	}
}
