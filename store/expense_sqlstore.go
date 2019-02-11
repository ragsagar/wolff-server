package store

import (
	"log"
	"net/url"
	"time"

	"github.com/go-pg/pg/orm"
	"github.com/go-pg/pg/urlvalues"
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

func (ess ExpenseSQLStore) GetExpenses(userId string, filter ExpenseFilter) ([]model.Expense, error) {
	var expenses []model.Expense
	err := ess.sqlStore.db.Model(&expenses).Column("expense.*").Relation("Category").Relation("Account").Where("expense.user_id = ?", userId).Apply(filter.Filter).Select()
	if err != nil {
		return nil, err
	}
	return expenses, nil
}

func (ess ExpenseSQLStore) StoreAccount(expenseAccount model.ExpenseAccount) error {
	err := ess.sqlStore.db.Insert(&expenseAccount)
	return err
}

func (ess ExpenseSQLStore) GetExpenseAccounts(userId string) ([]model.ExpenseAccount, error) {
	var expenseAccounts []model.ExpenseAccount
	err := ess.sqlStore.db.Model(&expenseAccounts).Where("user_id = ?", userId).Select()
	if err != nil {
		log.Println("Error in fetching accounts for user_id ", userId)
		return nil, err
	}
	return expenseAccounts, nil
}

func (ess ExpenseSQLStore) GetAccountByID(id string) (*model.ExpenseAccount, error) {
	expenseAccount := new(model.ExpenseAccount)
	err := ess.sqlStore.db.Model(expenseAccount).Where("expense_account.id = ?", id).Select()
	if err != nil {
		log.Println("Error in deleting: ", err.Error())
		return nil, err
	}
	return expenseAccount, nil
}

func (ess ExpenseSQLStore) DeleteAccount(expenseAccount *model.ExpenseAccount) error {
	return ess.sqlStore.db.Delete(expenseAccount)
}

// func (ess ExpenseSQLStore) DeleteExpenseWithUserID(id, userId string) error {
// 	return ess.sqlStore.db.Delete()
// }

type ExpenseFilter struct {
	*urlvalues.Pager
	// Filter(*orm.Query) (*orm.Query, error)
	year  int
	month int
}

func (f ExpenseFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.year > 0 {
		q = q.Where("EXTRACT(YEAR FROM expense.date) = ?", f.year)
	}

	if f.month > 0 {
		q = q.Where("EXTRACT(MONTH FROM expense.date) = ?", f.month)
	}

	q = q.Apply(f.Pager.Pagination)
	return q, nil
}

func (f *ExpenseFilter) ParseURLValues(values url.Values) {
	v := urlvalues.Values(values)
	currentTime := time.Now()

	year, err := v.Int("year")
	if err != nil || year == 0 {
		// default year to current year
		f.year = currentTime.Year()
	} else {
		f.year = year
	}

	month, err := v.Int("month")
	if err != nil || month == 0 {
		f.month = int(currentTime.Month())
	} else {
		f.month = month
	}

	f.Pager = urlvalues.NewPager(v)
}
