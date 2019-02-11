package store

import (
	"github.com/ragsagar/wolff/model"
)

// Store : Interface for the different store options
type Store interface {
	User() UserStore
	AuthToken() AuthTokenStore
	Expense() ExpenseStore
}

// UserStore : Interface for User store.
type UserStore interface {
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	StoreUser(user model.User) error
	UpdateUser(user model.User) error
}

// AuthTokenStore is an interface for AuthToken implementations.
type AuthTokenStore interface {
	Create(user *model.User) (*model.AuthToken, error)
	Find(token string) (*model.AuthToken, error)
}

// ExpenseStore is the interface that defines methods expected in ExpenseStore implemntations
type ExpenseStore interface {
	Store(expense *model.Expense) error
	GetByID(id string) (*model.Expense, error)
	GetExpenses(userId string, filter ExpenseFilter) ([]model.Expense, error)
	StoreAccount(expense model.ExpenseAccount) error
	GetExpenseAccounts(userId string) ([]model.ExpenseAccount, error)
	GetAccountByID(id string) (*model.ExpenseAccount, error)
	DeleteAccount(*model.ExpenseAccount) error
}
