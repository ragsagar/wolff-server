package store

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ragsagar/wolff/model"
	"github.com/stretchr/testify/suite"
)

type AuthTokenSQLStoreSuite struct {
	suite.Suite
	store SQLStore
	db    *sql.DB
}

func (s *AuthTokenSQLStoreSuite) SetupSuite() {
	s.T().Log("SetupSuite AuthToken running")
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
		`DROP TABLE auth_tokens`,
		`DROP TABLE users`,
	}
	for _, query := range queries {
		_, err := s.db.Query(query)
		if err != nil {
			s.T().Fatal(err)
		}
	}

	s.store = NewSQLStore(user, password, dbname, "localhost:5432")
	// query := `CREATE TABLE IF NOT EXISTS auth_tokens (
	//   key varchar(100) NOT NULL,
	//   user_id uuid,
	//   active boolean,
	//   expiry timestamptz,
	//   created_at timestamptz,
	//   updated_at timestamptz
	//   )`
	// _, err = s.db.Query(query)
	// if err != nil {
	// 	s.T().Fatal(err)
	// }
	// queries := []string{
	// 	`CREATE TABLE IF NOT EXISTS users (
	//   id uuid NOT NULL,
	//   email varchar(200) NOT NULL,
	//   password varchar(100) NOT NULL,
	//   name varchar(240),
	//   active boolean,
	//   created_at timestamptz,
	//   updated_at timestamptz
	//   )`,
	// 	`CREATE TABLE IF NOT EXISTS auth_tokens (
	//   key varchar(100) NOT NULL,
	//   user_id uuid,
	//   active boolean,
	//   expiry timestamptz
	//   )`,
	// }
	// for _, query := range queries {
	// 	_, err = s.db.Query(query)
	// 	if err != nil {
	// 		s.T().Fatal(err)
	// 	}
	// }

}

func (s *AuthTokenSQLStoreSuite) SetupTest() {
	s.T().Log("SetupTest AuthToken running")
	queries := []string{
		`TRUNCATE users`,
		`TRUNCATE auth_tokens`,
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

	testTokens := []struct {
		key    string
		expiry time.Time
		active bool
		userID string
	}{
		{"1234", time.Now().Add(time.Hour * 1), true, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"5678", time.Now().Add(time.Hour * 2), true, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
		{"2345", time.Now().Add(time.Hour * -3), true, "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"},
	}
	for _, t := range testTokens {
		_, err := s.db.Query("INSERT INTO auth_tokens (key, expiry, active, user_id) VALUES ($1, $2, $3, $4)", t.key, t.expiry, t.active, t.userID)
		if err != nil {
			log.Println("Error in inserting.")
			s.T().Fatal(err)
		}
	}
}

func TestAuthTokenSQLStoreSuite(t *testing.T) {
	s := new(AuthTokenSQLStoreSuite)
	suite.Run(t, s)
}

func (s *AuthTokenSQLStoreSuite) TestFind() {
	tokenString := "1234"
	aToken, err := s.store.AuthToken().Find(tokenString)
	if err != nil {
		s.T().Fatal(err)
	}
	if aToken.UserID != "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00" {
		s.T().Errorf("Error auth token user id is not matching. %s", aToken.UserID)
	}

	if aToken.User == nil {
		s.T().Errorf("User object is nil in AuthToken.")
	}

	if aToken.Key != tokenString {
		s.T().Errorf("Token is not matching.")
	}

	// Expired token
	tokenString = "2345"
	aToken, err = s.store.AuthToken().Find(tokenString)
	if aToken != nil && err == nil {
		s.T().Errorf("Find shouldn't return expired token. Error %s", err)
	}

	tokenString = "invalidtoken"
	aToken, err = s.store.AuthToken().Find(tokenString)
	if aToken != nil && err == nil {
		s.T().Errorf("Invalid token shouldn't return valid authToken. Error %s", err)
	}
}

func (s *AuthTokenSQLStoreSuite) TestCreate() {
	uid := "6a2191d9-77b2-4f14-b1c4-edb094be7cca"
	user := &model.User{ID: uid}
	authToken, err := s.store.AuthToken().Create(user)
	if err != nil {
		s.T().Fatal(err)
	}

	if !authToken.Expiry.After(time.Now()) {
		s.T().Errorf("New token shouldn't be already expired.")
	}
}
