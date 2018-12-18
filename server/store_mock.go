package server

import (
	"github.com/ragsagar/wolff/model"
	"github.com/ragsagar/wolff/store"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	userStore    *MockUserStore
	tokenStore   *MockAuthTokenStore
	expenseStore *MockExpenseStore
}

func NewMockStore() *MockStore {
	return &MockStore{
		userStore:    new(MockUserStore),
		tokenStore:   new(MockAuthTokenStore),
		expenseStore: new(MockExpenseStore),
	}
}

func (m MockStore) User() store.UserStore {
	return m.userStore
}

func (m MockStore) AuthToken() store.AuthTokenStore {
	return m.tokenStore
}

func (m MockStore) Expense() store.ExpenseStore {
	return m.expenseStore
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) GetUserByID(id string) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) GetUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserStore) StoreUser(user model.User) error {
	//args := m.Called(user)
	//return args.Error(0)
	return nil
}

func (m *MockUserStore) UpdateUser(user model.User) error {
	return nil
}

type MockAuthTokenStore struct {
	mock.Mock
}

func (m *MockAuthTokenStore) Find(token string) (*model.AuthToken, error) {
	args := m.Called(token)
	return args.Get(0).(*model.AuthToken), args.Error(1)
}

func (m *MockAuthTokenStore) Create(user *model.User) (*model.AuthToken, error) {
	args := m.Called(user)
	token := &model.AuthToken{UserID: user.ID, User: user, Key: "1234"}
	return token, args.Error(0)
}

type MockExpenseStore struct {
	mock.Mock
}

func (m MockExpenseStore) GetByID(id string) (*model.Expense, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Expense), args.Error(1)
}

func (m MockExpenseStore) Store(expense *model.Expense) error {
	expense.PreSave()
	return nil
}

func (m MockExpenseStore) GetExpenses(userId string) ([]model.Expense, error) {
	return nil, nil
}

func (m MockExpenseStore) StoreAccount(expenseAccount model.ExpenseAccount) error {
	expenseAccount.PreSave()
	return nil
}

func (m MockExpenseStore) GetExpenseAccounts(userId string) ([]model.ExpenseAccount, error) {
	// args := m.Called(userId)
	return nil, nil
}

func (m MockExpenseStore) GetAccountByID(id string) (*model.ExpenseAccount, error) {
	return nil, nil
}

func (m MockExpenseStore) DeleteAccount(expenseAccount *model.ExpenseAccount) error {
	return nil
}
