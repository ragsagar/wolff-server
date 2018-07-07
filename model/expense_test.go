package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpense(t *testing.T) {
	eAccount := &ExpenseAccount{
		ID:   "111",
		Name: "Grocery",
	}
	eCategory := &ExpenseCategory{
		ID:   "121",
		Name: "Category1",
	}
	expense := Expense{
		Amount:     100,
		AccountID:  eAccount.ID,
		Account:    eAccount,
		Category:   eCategory,
		CategoryID: eCategory.ID,
	}
	expense.PreSave()
	assert.NotNil(t, expense.ID, "Expense id is nil")

	expectedJSON, err := json.Marshal(expense)
	if err != nil {
		t.Fatal(err)
	}
	expenseJSON, err := expense.ToJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(expectedJSON) != string(expenseJSON) {
		t.Errorf("ToJSON response is not expected.")
	}
	expected := fmt.Sprintf("Expense<%s>", expense.ID)
	assert.Equal(t, expected, fmt.Sprintf("%s", expense))
}
