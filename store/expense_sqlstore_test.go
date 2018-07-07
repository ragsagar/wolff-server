package store

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/ragsagar/wolff/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExpenseSQLStoreSuite struct {
	suite.Suite
	store SQLStore
	db    *sql.DB
}

func (s *ExpenseSQLStoreSuite) SetupSuite() {
	s.T().Log("SetupSuite ExpenseSQLStore running")
	dbname := "wolffdb_test"
	user := "wolffuser"
	password := "password"
	connString := fmt.Sprintf("host=localhost port=5432 user=%s "+
		"password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db
	queries := []string{
		`DROP TABLE users`,
		`DROP TABLE expenses`,
		`DROP TABLE expense_categories`,
		`DROP TABLE expense_accounts`,
	}
	for _, query := range queries {
		_, err := s.db.Query(query)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	s.store = NewSQLStore(user, password, dbname, "localhost:5432")
}

func (s *ExpenseSQLStoreSuite) SetupTest() {
	s.T().Log("SetupTest ExpenseSQLStore running.")
	queries := []string{
		`TRUNCATE users`,
		`TRUNCATE expenses`,
		`TRUNCATE expense_categories`,
		`TRUNCATE expense_accounts`,
	}
	for _, query := range queries {
		_, err := s.db.Query(query)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	testUsers := []struct {
		uid      string
		email    string
		password string
		active   bool
	}{
		{"5d6e34c8-46b7-11e6-ba7c-cafec0ffee00", "testuser1@gmail.com", "password", true},
		{"6d6e34c8-56b7-11e6-ba7c-cafec0ffee00", "testuser2@gmail.com", "password", true},
	}

	for _, u := range testUsers {
		now := time.Now()
		hashedPassword, _ := model.HashPassword(u.password)
		_, err := s.db.Query("INSERT INTO users (id, email, password, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $5)", u.uid, u.email, hashedPassword, u.active, now)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	testAccounts := []struct {
		name    string
		id      string
		user_id string
	}{
		{"Grocery", "1234", "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"Restaurant", "2234", "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
	}

	for _, ta := range testAccounts {
		now := time.Now()
		_, err := s.db.Query("INSERT INTO expense_accounts (id, name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $4)", ta.id, ta.name, ta.user_id, now)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	testCategories := []struct {
		name   string
		id     string
		userID string
	}{
		{"Category1", "15678", "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"Category2", "25678", "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
	}

	for _, tc := range testCategories {
		now := time.Now()
		_, err := s.db.Query("INSERT INTO expense_categories (id, name, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $4)", tc.id, tc.name, tc.userID, now)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	testExpenses := []struct {
		id         string
		accountID  string
		categoryID string
		amount     float64
		userID     string
	}{
		{"14566", "1234", "15678", 50.0, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"24566", "1234", "15678", 60.0, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"34566", "1234", "25678", 70.0, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"44566", "2234", "25678", 100.0, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"54566", "2234", "15678", 60.0, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
	}
	for _, te := range testExpenses {
		now := time.Now()
		_, err := s.db.Query("INSERT INTO expenses (id, account_id, category_id, amount, user_id, created_at, updated_at, date) VALUES ($1, $2, $3, $4, $5, $6, $6, $6)",
			te.id, te.accountID, te.categoryID, te.amount, te.userID, now)
		if err != nil {
			s.T().Fatal(err)
		}
	}

}

func TestExpenseSQLStoreSuite(t *testing.T) {
	s := new(ExpenseSQLStoreSuite)
	suite.Run(t, s)
}

func (s *ExpenseSQLStoreSuite) TestGetByID() {
	eid := "14566"
	exp, err := s.store.Expense().GetByID(eid)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), exp.ID, eid, "Return expense object doesn't match.")
	assert.Equal(s.T(), exp.AccountID, "1234", "Expense account id doesn't match.")
	assert.Equal(s.T(), exp.CategoryID, "15678", "Expense category id doesn't match.")
	assert.NotNil(s.T(), exp.Account, "Expense account is nil after fetching from db.")
	assert.NotNil(s.T(), exp.Category, "Expense category is nil after fetching from db.")
	assert.NotNil(s.T(), exp.User, "User is nil after fetching from db.")
	assert.Equal(s.T(), "1234", exp.Account.ID, "Account id is not matching.")
	assert.Equal(s.T(), "15678", exp.Category.ID, "Category id is not matching.")
	assert.Equal(s.T(), "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00", exp.User.ID, "User id is not matching.")
	assert.Equal(s.T(), "Grocery", exp.Account.Name, "Account name is not matching.")
	assert.Equal(s.T(), "Category1", exp.Category.Name, "Category name is not matching.")
	assert.Equal(s.T(), exp.User.Email, "testuser1@gmail.com")

}

func (s *ExpenseSQLStoreSuite) TestStore() {
	expense := &model.Expense{
		AccountID:  "1234",
		CategoryID: "25678",
		Amount:     200,
		UserID:     "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00",
	}

	err := s.store.Expense().Store(expense)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), expense.ID, "Expense id is nil.")

	// Fetch the created account.
	exp, err := s.store.Expense().GetByID(expense.ID)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), exp.Account, "Expense account is nil after fetching from db.")
	assert.Equal(s.T(), expense.AccountID, exp.Account.ID)
	assert.Equal(s.T(), "Grocery", exp.Account.Name)
	assert.NotNil(s.T(), exp.Category, "Expense category is nil after fetching from db.")
	assert.Equal(s.T(), expense.CategoryID, exp.Category.ID)
	assert.Equal(s.T(), "Category2", exp.Category.Name)
}
