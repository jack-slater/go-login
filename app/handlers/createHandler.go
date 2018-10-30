package handlers

import (
	"net/http"
	"github.com/jack-slater/go-login/app/datastore"
	"github.com/jack-slater/go-login/app/model"
	"encoding/json"
	"strings"
	"fmt"
)


func CreateHandler(db datastore.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			saveUser(w, r, db)
		}
	})
}

func saveUser(w http.ResponseWriter, r *http.Request, db datastore.Database) {

	userDto := UserDTO{}
	if err := decode(r, &userDto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := model.NewUser(userDto.FirstName, userDto.LastName, userDto.Email, userDto.Login, userDto.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	saveToPostgres(db, user, w)
}

func saveToPostgres(db datastore.Database, user *model.User, w http.ResponseWriter) {
	createdUser, err := db.CreateUser(user)
	if err != nil {
		evaluateCreateErr(err, w)
		return
	}
	RespondWithJson(w, http.StatusCreated, createdUser)
}

func evaluateCreateErr(err error, w http.ResponseWriter) {
	var s string
	if strings.Contains(err.Error(), "email") {
		s = "email"
	} else {
		s = "login"
	}
	RespondWithError(w, http.StatusConflict, fmt.Sprintf("Could not create user - %s is not unique", s))
}

func decode(r *http.Request, u *UserDTO) error {
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		return err
	}
	return u.valid()
}