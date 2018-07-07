package store

import (
	"log"

	"github.com/ragsagar/wolff/model"
)

// ExpenseSQLStore is the SQL implementation of ExpenseStore interface.
type ExpenseSQLStore struct {
	sqlStore *SQLStore
}

// NewExpenseSQLStore returns new ExpenseSQLStore object.
func NewExpenseSQLStore(sqlStore SQLStore) *ExpenseSQLStore {
	return &ExpenseSQLStore{sqlStore: &sqlStore}
}

// Store saves the given expense object into database after populating ID, CreatedAt and UpdatedAt fields.
func (ess ExpenseSQLStore) Store(expense *model.Expense) error {
	expense.PreSave()
	err := ess.sqlStore.db.Insert(expense)
	return err
}

// GetByID fetches the Expense object with given id with its related ExpenseAccount, ExpenseCategory and User
func (ess ExpenseSQLStore) GetByID(id string) (*model.Expense, error) {
	expense := new(model.Expense)
	err := ess.sqlStore.db.Model(expense).Relation("Account").Relation("Category").Relation("User").Where("expense.id = ?", id).Select()
	if err != nil {
		log.Println("Error in fetching expense with id ", id)
		return nil, err
	}
	return expense, nil
}
