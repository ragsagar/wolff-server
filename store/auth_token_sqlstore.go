package store

import (
	"errors"
	"log"

	"github.com/ragsagar/wolff/model"
)

// AuthTokenSQLStore is the SQL implementation of AuthTokenStore interface
type AuthTokenSQLStore struct {
	sqlStore *SQLStore
}

// NewAuthTokenSQLStore returns new AuthTokenSQLStore object.
func NewAuthTokenSQLStore(sqlStore SQLStore) *AuthTokenSQLStore {
	return &AuthTokenSQLStore{sqlStore: &sqlStore}
}

// Create will create a new token and return it.
func (authTokenSQLStore AuthTokenSQLStore) Create(user *model.User) (*model.AuthToken, error) {
	authToken := model.NewAuthToken(user)
	authToken.PreSave()
	err := authTokenSQLStore.sqlStore.db.Insert(authToken)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

// Find returns the AuthToken instance with the given token key in database.
func (authTokenSQLStore AuthTokenSQLStore) Find(token string) (*model.AuthToken, error) {
	authToken := new(model.AuthToken)
	err := authTokenSQLStore.sqlStore.db.Model(authToken).Column("auth_token.*", "User").Relation("User").Where("auth_token.key = ? AND auth_token.expiry > NOW()", token).Select()
	if err != nil {
		log.Println("Error in fetching token: ", err.Error())
		return nil, errors.New("Token not found in db.")
	}
	return authToken, nil
}
