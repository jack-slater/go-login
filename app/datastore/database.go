package datastore

import "github.com/jack-slater/go-login/app/model"

type Database interface {
	GetUser(login, password string) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	Close()
}
