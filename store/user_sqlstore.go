package store

import "github.com/ragsagar/wolff/model"

// UserSQLStore : SQL implementation for UserStore interface.
type UserSQLStore struct {
	sqlStore *SQLStore
}

// NewUserSQLStore : Create new UserSQLStore with a SQLStore object.
func NewUserSQLStore(sqlStore SQLStore) *UserSQLStore {
	return &UserSQLStore{sqlStore: &sqlStore}
}

// GetUserByID : Fetch the user with given id as model.User
func (uss UserSQLStore) GetUserByID(id string) (*model.User, error) {
	user := model.User{ID: id}
	err := uss.sqlStore.db.Select(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail : Fetch the user with give email as model.User
func (uss UserSQLStore) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := uss.sqlStore.db.Model(&user).Where("email = ?", email).Select()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// StoreUser : Store the given model.User object to database.
func (uss UserSQLStore) StoreUser(user model.User) error {
	err := uss.sqlStore.db.Insert(&user)
	return err
}

// UpdateUser : Update the given model.User object in db.
func (uss UserSQLStore) UpdateUser(user model.User) error {
	err := uss.sqlStore.db.Update(&user)
	return err
}
