package model

import "time"

// AuthToken represent each token used for user authentication.
type AuthToken struct {
	Key    string
	User   *User
	UserID string
	Expiry time.Time
	Active bool
}

// NewAuthToken returns new AuthToken object.
func NewAuthToken(user *User) *AuthToken {
	return &AuthToken{Active: true, User: user, UserID: user.ID}
}

// PreSave populates the required fields before saving.
func (authToken *AuthToken) PreSave() {
	authToken.Key = GenerateTokenKey()
	authToken.Expiry = time.Now().Add(time.Hour * 72)
}
