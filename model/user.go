package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// User model
type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name,omitempty"`
	Active    bool      `json:"active"`
}

// SetPassword : Set new password for the user.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}

func (u User) String() string {
	return fmt.Sprintf("User<%s>", u.Email)
}

// ToJSON : Return the user object as json.
func (u User) ToJSON() ([]byte, error) {
	data, err := json.Marshal(u)
	return data, err
}

// PreSave : Call this method before passing the object to data store for saving.
// Populates the Id and time fields.
func (u *User) PreSave() {
	if u.ID == "" {
		u.ID = GenerateUUID()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
}
