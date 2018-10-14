package handlers

import (
	"net/http"
	"github.com/jack-slater/go-login/app/datastore"
	"github.com/jack-slater/go-login/app/model"
	"fmt"
	"encoding/json"
	"strings"
)

type UserDTO struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Login string `json:"login"`
	Password string `json:"password"`
}

func UserHandler(db datastore.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			saveUser(w, r, db)
		}

		if r.Method == "GET" {
			fmt.Fprintf(w, "wtf")
		}
	})
}

func saveUser(w http.ResponseWriter, r *http.Request, db datastore.Database) {

	decoder := json.NewDecoder(r.Body)
	userDto := UserDTO{}
	if err := decoder.Decode(&userDto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	user, err := model.NewUser(userDto.FirstName, userDto.LastName, userDto.Email, userDto.Login, userDto.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdUser, err := db.CreateUser(user)
	if err != nil {
		var errMsg = "is not unique"
		if strings.Contains(err.Error(), "email") {
			errMsg = "email " + errMsg
		} else {
			errMsg = "login " + errMsg
		}
		RespondWithError(w, http.StatusConflict, "Could not create user - " + errMsg)
	}

	RespondWithJson(w, 202, createdUser)
}