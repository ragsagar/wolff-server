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

func TestLoginUser(t *testing.T) {
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

	t1 := model.AuthToken{Key: "1234", UserID: expectedUser.ID, User: &expectedUser}

	testStore.tokenStore.On("Find", "1234").Return(&t1, nil)
	testStore.tokenStore.On("Create", mock.Anything).Return(nil)
	testStore.userStore.On("GetUserByEmail", expectedUser.Email).Return(&expectedUser, nil)

	body := map[string]string{}
	jsonData, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/api/users/login/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	srv := NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
	responseJSON, _ := ioutil.ReadAll(recorder.Body)
	expected := map[string]string{"error_message": "Email or password is missing."}
	expJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	if string(responseJSON) != string(expJSON) {
		t.Errorf("Returned error message not matching with expected. %s, %s", responseJSON, expJSON)
	}

	// reader := strings.NewReader("email=testuser5@gmail.com&password=password")
	body = map[string]string{"email": "testuser1@gmail.com", "password": "password"}
	jsonData, _ = json.Marshal(body)
	req, err = http.NewRequest("POST", "/api/users/login/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	srv.Routes.Root.ServeHTTP(recorder, req)
	var responseObj struct {
		AuthToken string `json:"auth_token"`
		UserID    string `json:"user_id"`
	}
	responseJSON, _ = ioutil.ReadAll(recorder.Body)
	json.Unmarshal(responseJSON, &responseObj)
	assert.Equal(t, "1234", responseObj.AuthToken, "Generated auth token not matching.")
	assert.Equal(t, expectedUser.ID, responseObj.UserID, "User Id in response is not what expected.")
	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

// func TestGetUserById(t *testing.T) {
// 	testStore := NewMockStore()
//
// 	expectedUser := model.User{
// 		Id:        "b89505a4-a451-45e5-912e-4ef8c1441be6",
// 		Email:     "testuser1@gmail.com",
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 		Active:    true,
// 	}
// 	t1 := model.AuthToken{Key: "1234", UserID: expectedUser.Id, User: &expectedUser}
//
// 	testStore.tokenStore.On("Find", "1234").Return(&t1, nil)
// 	testStore.userStore.On("GetUserByID", expectedUser.Id).Return(&expectedUser, nil)
// 	// oat := model.OauthAccessToken{Token: "1234", Scopes: "user,service"}
// 	// testStore.tokenStore.On("FindByToken", "1234").Return(&oat, nil)
//
// 	req, err := http.NewRequest("GET", "/api/users/b89505a4-a451-45e5-912e-4ef8c1441be6", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req.Header.Add("Authorization", "1234")
// 	recorder := httptest.NewRecorder()
// 	srv := NewServer(testStore)
// 	srv.Routes.Root.ServeHTTP(recorder, req)
// 	if status := recorder.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v", status)
// 	}
//
// 	expectedUserJson, _ := expectedUser.ToJSON()
// 	responseJson, _ := ioutil.ReadAll(recorder.Body)
// 	if string(responseJson) != string(expectedUserJson) {
// 		t.Errorf("Returned user not matching with expected user.")
// 	}
// 	testStore.userStore.AssertExpectations(t)
// }

func TestCreateUser(t *testing.T) {
	testStore := NewMockStore()
	id := model.GenerateUUID()
	expectedUser := model.User{
		ID:        id,
		Email:     "testuser2@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    true,
	}
	t1 := model.AuthToken{Key: "1234", UserID: expectedUser.ID, User: &expectedUser}

	testStore.tokenStore.On("Find", "1234").Return(&t1, nil)
	testStore.tokenStore.On("Create", mock.Anything).Return(nil)
	testStore.userStore.On("GetUserByID", expectedUser.ID).Return(&expectedUser, nil)
	testStore.userStore.On("StoreUser", expectedUser).Return(nil)
	// oat := model.OauthAccessToken{Token: "1234", Scopes: "user,service"}
	// testStore.tokenStore.On("FindByToken", "1234").Return(&oat, nil)

	body := map[string]string{}
	jsonData, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", "/api/users/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "1234")
	recorder := httptest.NewRecorder()
	srv := NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
	responseJSON, _ := ioutil.ReadAll(recorder.Body)
	expected := map[string]string{"error_message": "Email or password is missing."}
	expJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatal(err)
	}
	if string(responseJSON) != string(expJSON) {
		t.Errorf("Returned error message not matching with expected. %s, %s", responseJSON, expJSON)
	}

	// reader := strings.NewReader("email=testuser5@gmail.com&password=password")
	body = map[string]string{"email": "testuser@gmail.com", "password": "password"}
	jsonData, _ = json.Marshal(body)

	req, err = http.NewRequest("POST", "/api/users/", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "1234")
	req.Header.Add("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	//srv = NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	log.Println(recorder.Body)
	var responseObj struct {
		AuthToken string `json:"auth_token"`
		UserId    string `json:"user_id"`
	}
	responseJSON, _ = ioutil.ReadAll(recorder.Body)
	json.Unmarshal(responseJSON, &responseObj)
	assert.Equal(t, "1234", responseObj.AuthToken, "Generated auth token not matching.")
	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v", status)
	}
}

func TestGetUserProfile(t *testing.T) {
	testStore := NewMockStore()

	expectedUser := model.User{
		ID:        "b89505a4-a451-45e5-912e-4ef8c1441be6",
		Email:     "testuser1@gmail.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Active:    true,
	}
	t1 := model.AuthToken{Key: "1234", UserID: expectedUser.ID, User: &expectedUser}

	testStore.tokenStore.On("Find", "1234").Return(&t1, nil)
	// testStore.userStore.On("GetUserByID", expectedUser.Id).Return(&expectedUser, nil)
	// oat := model.OauthAccessToken{Token: "1234", Scopes: "user,service"}
	// testStore.tokenStore.On("FindByToken", "1234").Return(&oat, nil)

	req, err := http.NewRequest("GET", "/api/users/profile/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "1234")
	recorder := httptest.NewRecorder()
	srv := NewServer(testStore)
	srv.Routes.Root.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	expectedUserJSON, _ := expectedUser.ToJSON()
	responseJSON, _ := ioutil.ReadAll(recorder.Body)
	if string(responseJSON) != string(expectedUserJSON) {
		t.Errorf("Returned user not matching with expected user.")
	}
	testStore.userStore.AssertExpectations(t)
}
