package server

import (
	"errors"
	"testing"

	"github.com/ragsagar/wolff/model"
)

func TestValidateToken(t *testing.T) {
	testStore := NewMockStore()
	// oat1 := model.OauthAccessToken{Token: "1234", Scopes: "user,service"}
	// oat2 := model.OauthAccessToken{Token: "5678", Scopes: "user"}
	t1 := model.AuthToken{Key: "1234"}
	t2 := model.AuthToken{Key: "5678"}
	testStore.tokenStore.On("Find", "1234").Return(&t1, nil)
	testStore.tokenStore.On("Find", "5678").Return(&t2, nil)
	testStore.tokenStore.On("Find", "1111").Return(&model.AuthToken{}, errors.New("Token not found in db"))
	testStore.tokenStore.On("Find", "").Return(nil, errors.New("Required service token"))
	// // Giving &oat1 as dummy since we are not able to pass nil there.
	// testStore.tokenStore.On("FindByToken", "1111").Return(&oat1, errors.New("Doesn't exist in db."))

	_, err := validateToken(testStore, "1234")
	if err != nil {
		t.Fatal(err)
	}

	_, err = validateToken(testStore, "5678")
	if err != nil {
		t.Fatal(err)
	}

	_, err = validateToken(testStore, "")
	if err.Error() != "Token missing" {
		t.Fatal(err)
	}

	// Non existing token
	_, err = validateToken(testStore, "1111")
	if err.Error() != "Token not found" {
		t.Fatal(err)
	}

}
