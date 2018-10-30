package model

import (
	"github.com/jack-slater/go-login/app/helpers"
	"errors"
)

//type UserSchema interface {
//	getUser(login, password string) (*model.User, error)
//	createUser(user model.User) (*model.User, error)
//}

type User struct {
	Id int `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Email string `json:"Email"`
	Login string `json:"-"`
	Password string `json:"-"`
}

func NewUser(firstName, lastName, email, login, password string) (*User, error) {

	if !helpers.VerifyEmailFormat(email) {
		return nil, errors.New("email is incorrect format. Cannot create user")
	}

	hashedPassword := helpers.IncryptPassword(password)
	return &User{ FirstName:firstName, LastName:lastName, Email:email, Login:login, Password:hashedPassword }, nil
}