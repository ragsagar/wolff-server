package store

import "github.com/go-pg/pg"

// SQLStore : SQL implementation of store.
type SQLStore struct {
	userSQLStore *UserSQLStore
	db           *pg.DB
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
	return sqlStore
}

// User : Return UserSQLStore, required to implement Store interface.
func (sqlStore SQLStore) User() UserStore {
	return sqlStore.userSQLStore
}
