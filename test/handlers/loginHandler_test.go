package handlers_test

import (
	"github.com/jack-slater/go-login/app/datastore"
	"github.com/jack-slater/go-login/app/model"
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/jack-slater/go-login/app/handlers"
	"bytes"
	"github.com/jack-slater/go-login/app/helpers"
	"fmt"
	"errors"
)

const (
	PASSWORD = "password"
	EMAIL    = "jack@email.com"
)

var hashedPassword = helpers.IncryptPassword(PASSWORD)
var user = model.User{1, "Jack", "Slater", EMAIL,
	"slater100", hashedPassword}

type MockLoginDb struct{ *datastore.PostgresDatastore }

func (db *MockLoginDb) GetUser(login, password string) (*model.User, error) {
	e := errors.New("unauthorised")
	if login != user.Email {
		return nil, e
	}

	if helpers.IncryptPassword(password) != hashedPassword {
		return nil, e
	}

	return &user, nil
}

func (db *MockLoginDb) Close()                                           {}
func (db *MockLoginDb) CreateUser(user *model.User) (*model.User, error) { return nil, nil }

func TestLoginHandler_success(t *testing.T) {
	payload := []byte(fmt.Sprintf(`{"email": "%v", "password": "%v"}`, EMAIL, PASSWORD))

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	rr := executeLoginResponse(req)

	if statusCode := rr.Code; statusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			statusCode, http.StatusOK)
	}

	want := `{"ID":1,"FirstName":"Jack","LastName":"Slater","Email":"jack@email.com"}`
	if got := rr.Body.String(); got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func TestLoginHandler_unknown(t *testing.T) {
	payload := []byte(fmt.Sprintf(`{"email": "%v", "password": "%v"}`, "unknown@email.com", PASSWORD))

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	rr := executeLoginResponse(req)

	if statusCode := rr.Code; statusCode != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			statusCode, http.StatusUnauthorized)
	}

	want := `{"error":"Incorrect login credentials"}`
	if got := rr.Body.String(); got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func TestLoginHandler_unauthorised(t *testing.T) {
	payload := []byte(fmt.Sprintf(`{"email": "%v", "password": "%v"}`, EMAIL, "incorrectPassword"))

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	rr := executeLoginResponse(req)

	if statusCode := rr.Code; statusCode != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			statusCode, http.StatusUnauthorized)
	}

	want := `{"error":"Incorrect login credentials"}`
	if got := rr.Body.String(); got != want {
		t.Errorf("handler returned unexpected body: got %v want %v", got, want)
	}
}

func TestLoginHandler_payloadIncorrectFormat(t *testing.T) {
	payloadIncorrectJson := []byte(fmt.Sprintf(`{email: "%v"}`, EMAIL))

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payloadIncorrectJson))
	rr := executeLoginResponse(req)

	if statusCode := rr.Code; statusCode != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			statusCode, http.StatusBadRequest)
	}
}

func executeLoginResponse(req *http.Request) *httptest.ResponseRecorder {
	mockDb := MockLoginDb{}
	rr := httptest.NewRecorder()
	handler := handlers.LoginHandler(&mockDb)
	handler.ServeHTTP(rr, req)
	return rr
}
