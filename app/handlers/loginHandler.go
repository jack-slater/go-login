package handlers

import (
	"github.com/jack-slater/go-login/app/datastore"
	"net/http"
	"encoding/json"
	"log"
)

type LoginDTO struct {
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
}

func LoginHandler(db datastore.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			loginUser(w, r, db)
		}
	})
}

func loginUser(w http.ResponseWriter, r *http.Request, db datastore.Database) {

	loginDto := LoginDTO{}
	if err := json.NewDecoder(r.Body).Decode(&loginDto); err != nil {
		log.Print(err)
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	user, err := db.GetUser(loginDto.Email, loginDto.Password)

	if err != nil {
		log.Print(err)
		RespondWithError(w, http.StatusUnauthorized, "Incorrect login credentials")
		return
	}

	RespondWithJson(w, http.StatusOK, user)
}