package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ragsagar/wolff/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMockStoreData(t *testing.T) *MockStore {
	testStore := NewMockStore()

	hashedPassword, err := model.HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	expectedUser := model.User{
		ID:        "b89505a4-a451-45e5-912e-4ef8c1441be6",
		Email:     "testuser1@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  hashedPassword,
		Active:    true,
	}

	eAccount := &model.ExpenseAccount{
		ID:   "111",
		Name: "Grocery",
	}
	eCategory := &model.ExpenseCategory{
		ID:   "121",
		Name: "Category1",
	}
	expense := model.Expense{
		Amount:     100,
		AccountID:  eAccount.ID,
		Account:    eAccount,
		Category:   eCategory,
		CategoryID: eCategory.ID,
		ID:         "9123",
	}

	t1 := model.AuthToken{Key: "1234", UserID: expectedUser.ID, User: &expectedUser}

	testStore.tokenStore.On("Find", "1234").Return(&t1, nil)
	testStore.tokenStore.On("Create", mock.Anything).Return(nil)
	testStore.userStore.On("GetUserByEmail", expectedUser.Email).Return(&expectedUser, nil)
	testStore.expenseStore.On("GetByID", "9123").Return(expense, nil)
	testStore.userStore.On("Store", mock.Anything).Return(nil)
	return testStore
}

func TestGetExpense(t *testing.T) {
	testStore := setupMockStoreData(t)

	// Empty body. Check if all error messages are appearing.
	body := map[string]string{}
	jsonData, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/api/expenses/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	srv := NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
	req.Header.Add("Authorization", "1234")
	recorder = httptest.NewRecorder()
	srv = NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
	responseBody, _ := ioutil.ReadAll(recorder.Body)
	var expected struct {
		Errors map[string][]string
	}
	err = json.Unmarshal(responseBody, &expected)
	if err != nil {
		log.Println(err)
	}
	assert.Contains(t, expected.Errors, "amount")
	assert.Contains(t, expected.Errors, "title")
	assert.Contains(t, expected.Errors, "date")
	assert.Equal(t, expected.Errors["amount"], []string{"amount field is required."})
	assert.Equal(t, expected.Errors["title"], []string{"title field is required."})
	assert.Equal(t, expected.Errors["date"], []string{"date field is required."})
	// if err := json.NewDecoder(recorder.Body).Decode(&responseJSON); err != nil {
	// 	t.Fatal(err)
	// }
	// log.Println(responseJSON)
	reqBody := map[string]interface{}{
		"amount":      100.0,
		"title":       "Nesto",
		"category_id": "121",
		"account_id":  "111",
		"date":        "2006-01-02T15:04:05Z",
	}
	jsonData, _ = json.Marshal(reqBody)
	req, err = http.NewRequest("POST", "/api/expenses/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "1234")
	recorder = httptest.NewRecorder()
	srv = NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
	responseBody, _ = ioutil.ReadAll(recorder.Body)
	expense := model.Expense{}
	if err := json.Unmarshal(responseBody, &expense); err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, expense.ID)
	assert.Equal(t, "111", expense.AccountID)
	assert.Equal(t, "b89505a4-a451-45e5-912e-4ef8c1441be6", expense.UserID)
	assert.Equal(t, "121", expense.CategoryID)
	assert.Equal(t, 100.0, expense.Amount)
}
