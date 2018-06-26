package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	Id            string    `json:"id"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	Email         string    `json:"email"`
	Password      string    `json:"-"`
	EmailVerified bool      `json:"-"`
	Active        bool      `json:"active"`
}

func (u *User) SetPassword(rawPassword string) {
	u.Password = HashPassword(rawPassword)
}

func (u User) String() string {
	return fmt.Sprintf("User<%s>", u.Email)
}

func (u User) ToJson() ([]byte, error) {
	data, err := json.Marshal(u)
	return data, err
}

func (u *User) PreSave() {
	if u.Id == "" {
		u.Id = GenerateUUID()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
}
