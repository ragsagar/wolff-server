package model

import (
	"time"
)

// Expense represent each expense data created by an user.
type Expense struct {
	ID         string
	Account    *ExpenseAccount
	AccountID  string
	Date       time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Category   *ExpenseCategory
	CategoryID string
	Amount     float64
	User       *User
	UserID     string
}

// ExpenseCategory for keeping each expense categories.
type ExpenseCategory struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
	UserID    string
}

// ExpenseAccount represent different expense accounts.
type ExpenseAccount struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      *User
	UserID    string
}

// PreSave popualtes ID, CreateAt and UpdatedAt fields. Call this before saving to db.
func (e *Expense) PreSave() {
	e.ID = GenerateUUID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}
