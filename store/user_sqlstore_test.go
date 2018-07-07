package store

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ragsagar/wolff/model"
	"github.com/stretchr/testify/suite"
)

type UserSQLStoreSuite struct {
	suite.Suite
	store SQLStore
	db    *sql.DB
}

func (s *UserSQLStoreSuite) SetupSuite() {
	s.T().Log("SetupSuite running for User.")
	dbname := "wolffdb_test"
	user := "wolffuser"
	password := "password"
	s.store = NewSQLStore(user, password, dbname, "localhost:5432")
	connString := fmt.Sprintf("host=localhost port=5432 user=%s "+
		"password=%s dbname=%s sslmode=disable", user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db

	queries := []string{
		`DROP TABLE users`,
		`CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL,
    email varchar(200) NOT NULL,
    password varchar(100) NOT NULL,
    name varchar(240),
    active boolean,
    created_at timestamptz,
    updated_at timestamptz
    )`,
	}
	for _, query := range queries {
		_, err = s.db.Query(query)
		if err != nil {
			s.T().Fatal(err)
		}
	}
}

func (s *UserSQLStoreSuite) SetupTest() {
	s.T().Log("SetupTest User running.")

	// Create some data for tests.
	_, err := s.db.Query("DELETE FROM users")
	if err != nil {
		s.T().Fatal(err)
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
	// uid := "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"
	// _, err = s.db.Query("INSERT INTO users (id, email, active) VALUES ($1, 'testuser1@gmail.com', true)", uid)
	// if err != nil {
	// 	s.T().Fatal(err)
	// }
	//
	// uid = "6d6e34c8-56b7-11e6-ba7c-cafec0ffee00"
	// _, err = s.db.Query("INSERT INTO users (id, email, active) VALUES ($1, 'testuser2@gmail.com', true)", uid)
	// if err != nil {
	// 	s.T().Fatal(err)
	// }
}

func TestUserSqlStoreSuite(t *testing.T) {
	s := new(UserSQLStoreSuite)
	suite.Run(t, s)

	// s2 := new(AuthTokenSQLStoreSuite)
	// suite.Run(t, s2)
}

func (s *UserSQLStoreSuite) TestGetUserById() {
	user, err := s.store.User().GetUserByID("5d6e34c8-46b7-11e6-ba7c-cafec0ffee00")
	if err != nil {
		s.T().Fatal(err)
	}
	if user.Email != "testuser1@gmail.com" {
		s.T().Errorf("User email id check failed.")
	}
}

func (s *UserSQLStoreSuite) TestGetUserByEmail() {
	user, err := s.store.User().GetUserByEmail("testuser1@gmail.com")
	if err != nil {
		s.T().Fatal(err)

	}

	if user.Email != "testuser1@gmail.com" {
		s.T().Errorf("Get user by email id returned wrong data")
	}

	if user.ID != "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00" {
		s.T().Errorf("Get user by email id returned wrong id.")
	}
}

func (s *UserSQLStoreSuite) TestStoreUser() {
	user := model.User{
		ID:    "5d6e34c8-46b7-11e6-ba7c-cafec0ffee12",
		Email: "testuser3@gmail.com",
	}
	user.SetPassword("password")
	err := s.store.User().StoreUser(user)
	if err != nil {
		s.T().Fatal(err)
	}

	res, err := s.db.Query("SELECT COUNT(*) FROM users WHERE email = $1", "testuser3@gmail.com")
	if err != nil {
		s.T().Fatal(err)
	}
	var count int
	for res.Next() {
		err = res.Scan(&count)
		if err != nil {
			s.T().Error(err)
		}
	}

	if count != 1 {
		s.T().Errorf("Incorrect count 1, StoreUser didn't work properly.")
	}

	res, err = s.db.Query("SELECT id, email, password from users WHERE email = $1", "testuser3@gmail.com")
	if err != nil {
		s.T().Fatal(err)
	}

	var user1 model.User
	for res.Next() {
		err := res.Scan(&user1.ID, &user1.Email, &user1.Password)
		if err != nil {
			s.T().Error(err)
		}
	}
	if user1 != user {
		s.T().Errorf("Data didnt get stored properly.")
	}

}

func (s *UserSQLStoreSuite) TestUpdateUser() {

	uid := "5d6e34c8-46b7-11e6-ba7c-cafec0ffee00"
	res, err := s.db.Query("SELECT id, email, password, active, created_at, updated_at from users WHERE id = $1", uid)
	if err != nil {
		s.T().Fatal(err)
	}

	var user1, user2 model.User
	for res.Next() {
		err = res.Scan(&user1.ID, &user1.Email, &user1.Password, &user1.Active, &user1.CreatedAt, &user1.UpdatedAt)
		if err != nil {
			s.T().Error(err)
		}
	}

	user1.Email = "someone1@gmail.com"
	err = s.store.User().UpdateUser(user1)
	if err != nil {
		s.T().Fatal(err)
	}

	res, err = s.db.Query("SELECT id, email, password, active, created_at, updated_at from users WHERE id = $1", uid)
	if err != nil {
		s.T().Fatal(err)
	}
	for res.Next() {
		err := res.Scan(&user2.ID, &user2.Email, &user2.Password, &user2.Active, &user2.CreatedAt, &user2.UpdatedAt)
		if err != nil {
			s.T().Error(err)
		}
	}

	if user2 != user1 {
		s.T().Errorf("Updating user failed.")
	}
}
