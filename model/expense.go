package model

import (
	"encoding/json"
	"fmt"
	"time"
)

// Expense represent each expense data created by an user.
type Expense struct {
	ID         string           `json:"id"`
	Account    *ExpenseAccount  `json:"account"`
	AccountID  string           `json:"account_id"`
	Date       time.Time        `json:"date"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Category   *ExpenseCategory `json:"category" pg:",fk:category_id"`
	CategoryID string           `json:"category_id"`
	Amount     float64          `json:"amount"`
	User       *User            `json:"user"`
	UserID     string           `json:"user_id"`
	Title      string           `json:"title"`
}

// String return the string representation of Expense object.
func (e Expense) String() string {
	return fmt.Sprintf("Expense<%s>", e.ID)
}

// PreSave popualtes ID, CreateAt and UpdatedAt fields. Call this before saving to db.
func (e *Expense) PreSave() {
	e.ID = GenerateUUID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

// ToJSON returns expense object as json
func (e Expense) ToJSON() ([]byte, error) {
	data, err := json.Marshal(e)
	return data, err
}

// ExpenseCategory for keeping each expense categories.
type ExpenseCategory struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user"`
	UserID    string    `json:"user_id"`
}

// ExpenseAccount represent different expense accounts.
type ExpenseAccount struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"-"`
	UserID    string    `json:"user_id"`
}

func (e ExpenseAccount) String() string {
	return fmt.Sprintf("ExpenseAccount<%s>", e.Name)
}

func (e *ExpenseAccount) PreSave() {
	e.ID = GenerateUUID()
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

func (e ExpenseAccount) ToJSON() ([]byte, error) {
	data, err := json.Marshal(e)
	return data, err
}
