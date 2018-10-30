package handlers_test

import (
	"net/http"
	"testing"

	"github.com/jack-slater/go-login/app/datastore"
	"github.com/jack-slater/go-login/app/handlers"
	"net/http/httptest"
	"github.com/jack-slater/go-login/app/model"
	"bytes"
	"github.com/pkg/errors"
)

type MockCreateDb struct{ *datastore.PostgresDatastore }

func (db *MockCreateDb) GetUser(login, password string) (*model.User, error) { return nil, nil }
func (db *MockCreateDb) Close()                               {}

func (db *MockCreateDb) CreateUser(user *model.User) (*model.User, error) {
	createdUser := model.User{2, "Jack", "Slater", "jack@slater.com", "slater001", "hashedPassowrd"}
	if createdUser.Email == user.Email {
		return nil, errors.New("email error")
	}

	if createdUser.Login == user.Login {
		return nil, errors.New("login error")
	}

	user.Id = 1
	return user, nil
}

func TestCreateHandler_successfulCreation(t *testing.T) {

	payload := []byte(`{"login": "jjj",
	"firstName": "Jim",
	"lastName": "Jime",
	"email": "time@gmail.com",
	"password": "sublime"}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(payload))
	rr := executeCreateResponse(req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	expected := `{"ID":1,"FirstName":"Jim","LastName":"Jime","Email":"time@gmail.com"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateHandler_alreadyUsedLogin(t *testing.T) {

	payload := []byte(`{"login": "slater001",
	"firstName": "Jim",
	"lastName": "Jime",
	"email": "time@gmail.com",
	"password": "sublime"}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(payload))
	rr := executeCreateResponse(req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	expected := `{"error":"Could not create user - login is not unique"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateHandler_alreadyUsedEmail(t *testing.T) {

	payload := []byte(`{"login": "jjj",
	"firstName": "Jim",
	"lastName": "Jime",
	"email": "jack@slater.com",
	"password": "sublime"}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(payload))
	rr := executeCreateResponse(req)

	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}

	expected := `{"error":"Could not create user - email is not unique"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateHandler_payloadIncorrectFormat(t *testing.T) {

	payloadWithNoLogin := []byte(`{
	"firstName": "Jim",
	"lastName": "Jime",
	"email": "jack@slater.com",
	"password": "sublime"}`)

	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(payloadWithNoLogin))
	rr := executeCreateResponse(req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := `{"error":"Invalid request payload"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func executeCreateResponse(req *http.Request) *httptest.ResponseRecorder {
	mockDb := MockCreateDb{}
	rr := httptest.NewRecorder()
	handler := handlers.CreateHandler(&mockDb)
	handler.ServeHTTP(rr, req)
	return rr
}
