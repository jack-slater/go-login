package datastore_test

import (
	_ "github.com/lib/pq"
	"testing"
	"fmt"
	"os"
	"github.com/jack-slater/go-login/app/datastore"
	"log"
	"github.com/jack-slater/go-login/app/model"
	"github.com/jack-slater/go-login/app/helpers"
	"reflect"
)

type TestPostgres struct {
	ds *datastore.PostgresDatastore
}

var testPostgres TestPostgres

func TestMain(m *testing.M) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=5433",
		"postgres-test", "password", "db_test")

	db, err := datastore.NewPostgresDataStore(connectionString)

	if err != nil {
		log.Print(err)
	}

	testPostgres.ds = db

	code := m.Run()
	clearTable()
	os.Exit(code)
}

func TestCreateUser_success(t *testing.T) {
	clearTable()

	firstName, lastName, email, login, password := "jack", "slater", "j@j.com", "slater100", "password"
	hashedPassword := helpers.IncryptPassword(password)

	user, _ := model.NewUser(firstName, lastName, email, login, password)

	want := model.User{Id: 1, FirstName: firstName, LastName: lastName, Email: email, Login: login, Password: hashedPassword}
	got, _ := testPostgres.ds.CreateUser(user)

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("got user: %v but want user: %v", got, &want)
	}

	secondUser, _ := model.NewUser(firstName, lastName, "new@email.com", "newLogin", password)

	wantSecond := model.User{Id: 2, FirstName: firstName, LastName: lastName, Email: "new@email.com", Login: "newLogin", Password: hashedPassword}
	gotSecond, _ := testPostgres.ds.CreateUser(secondUser)

	if !reflect.DeepEqual(gotSecond, &wantSecond) {
		t.Errorf("got user: %v but want user: %v", gotSecond, &wantSecond)
	}
}

func TestCreateUser_idIncrements(t *testing.T) {
	clearTable()

	firstName, lastName, email, login, password, secondEmail, secondLogin :=
		"jack", "slater", "j@j.com", "slater100", "password", "second@email.com", "secondLogin"

	firstUser, _ := model.NewUser(firstName, lastName, email, login, password)
	secondUser, _ := model.NewUser(firstName, lastName, secondEmail, secondLogin, password)

	gotFirst, _ := testPostgres.ds.CreateUser(firstUser)
	gotSecond, _ := testPostgres.ds.CreateUser(secondUser)

	if gotFirst.Id != 1 {
		t.Errorf("got first created user id: %v  but want user id: %v", gotFirst.Id, 1)
	}

	if gotSecond.Id != 2 {
		t.Errorf("got second created user id: %v  but want user id: %v", gotSecond.Id, 1)
	}
}

func TestCreateUser_emailNotUnique(t *testing.T) {
	clearTable()

	firstName, lastName, email, login, password := "jack", "slater", "j@j.com", "slater100", "password"

	user, _ := model.NewUser(firstName, lastName, email, login, password)
	testPostgres.ds.CreateUser(user)

	newUser, _ := model.NewUser("newName", "newName", email, "newLogin", "newPassword")
	_, err := testPostgres.ds.CreateUser(newUser)

	if err == nil {
		t.Errorf("expected error %v", err)
	}

}

func TestCreateUser_loginNotUnique(t *testing.T) {
	clearTable()

	firstName, lastName, email, login, password := "jack", "slater", "j@j.com", "slater100", "password"

	user, _ := model.NewUser(firstName, lastName, email, login, password)
	testPostgres.ds.CreateUser(user)

	newUser, _ := model.NewUser("newName", "newName", "new@email.com", login, "newPassword")
	_, err := testPostgres.ds.CreateUser(newUser)

	if err == nil {
		t.Errorf("expected error %v", err)
	}

}

func clearTable() {
	testPostgres.ds.DB.Exec(`DELETE FROM "user"`)
	testPostgres.ds.DB.Exec(`ALTER SEQUENCE user_id_seq RESTART WITH 1`)
}
