package store

import (
	"log"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/ragsagar/wolff/model"
)

// SQLStore : SQL implementation of store.
type SQLStore struct {
	userSQLStore   *UserSQLStore
	authTokenStore *AuthTokenSQLStore
	expenseStore   *ExpenseSQLStore
	db             *pg.DB
}

// NewSQLStore : Initializes SQLStore object, establishes db connection with
// the given parameters and return it.
func NewSQLStore(user, password, database, addr string) SQLStore {
	sqlStore := SQLStore{}
	sqlStore.db = pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Database: database,
		Addr:     addr,
	})
	sqlStore.userSQLStore = NewUserSQLStore(sqlStore)
	sqlStore.authTokenStore = NewAuthTokenSQLStore(sqlStore)
	sqlStore.expenseStore = NewExpenseSQLStore(sqlStore)
	createSchema(sqlStore.db)
	return sqlStore
}

// User : Return UserSQLStore, required to implement Store interface.
func (sqlStore SQLStore) User() UserStore {
	return sqlStore.userSQLStore
}

// AuthToken returns AuthTokenSQLStore which is required to implement Store interface.
func (sqlStore SQLStore) AuthToken() AuthTokenStore {
	return sqlStore.authTokenStore
}

// Expense returns ExpenseSQLStore to implement Store interface
func (sqlStore SQLStore) Expense() ExpenseStore {
	return sqlStore.expenseStore
}

func createSchema(db *pg.DB) {
	log.Println("Creating schema.")
	// queries := []string{
	// 	`CREATE TABLE IF NOT EXISTS users (
	//     id uuid NOT NULL,
	//     email varchar(200) NOT NULL,
	//     password varchar(100) NOT NULL,
	//     name varchar(240),
	//     active boolean,
	//     created_at timestamptz,
	//     updated_at timestamptz
	//     )`,
	// 	`CREATE TABLE IF NOT EXISTS auth_tokens (
	//     key varchar(100) NOT NULL,
	//     user_id uuid,
	//     active boolean,
	//     expiry timestamptz
	//     )`,
	// }
	// for _, query := range queries {
	// 	_, err := db.Query(query)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	models := []interface{}{
		(*model.User)(nil),
		(*model.AuthToken)(nil),
		(*model.Expense)(nil),
		(*model.ExpenseAccount)(nil),
		(*model.ExpenseCategory)(nil),
	}
	for _, model := range models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			panic(err)
		}
	}
}
