package model

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
	//TODO generate id, verify login is unique and hash password
	return &User{ FirstName:firstName, LastName:lastName, Email:email, Login:login, Password:password }, nil
}