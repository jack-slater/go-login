package model_test

import (
	"reflect"
	"testing"
	"github.com/jack-slater/go-login/app/helpers"
	"errors"
	"github.com/jack-slater/go-login/app/model"
)

func TestNewUser(t *testing.T) {

	firstName, lastName, email, login, password := "jack", "slater", "j@j.com", "slater100", "password"
	hashedPassword := helpers.IncryptPassword(password)

	type args struct {
		firstName string
		lastName  string
		email     string
		login     string
		password  string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.User
		wantErr error
	}{
		{
			"creates new user",
			args{firstName, lastName, email, login, password},
			&model.User{FirstName: firstName, LastName: lastName, Email: email, Login: login, Password: hashedPassword},
			nil},
		{
			"email error with incorrect email",
			args{firstName, lastName, "incorrectEmail", login, password},
			nil,
			errors.New("email is incorrect format. Cannot create user")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := model.NewUser(tt.args.firstName, tt.args.lastName, tt.args.email, tt.args.login, tt.args.password)
			if (err != nil) && err.Error() != tt.wantErr.Error() {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
