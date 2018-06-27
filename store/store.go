package store

import (
	"github.com/ragsagar/wolff/model"
)

// Store : Interface for the store.
type Store interface {
	User() UserStore
}

// UserStore : Interface for User store.
type UserStore interface {
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	StoreUser(user model.User) error
	UpdateUser(user model.User) error
}
