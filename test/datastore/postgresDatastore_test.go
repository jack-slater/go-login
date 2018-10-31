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

const (
	FIRST_NAME = "jack"
	LAST_NAME  = "slater"
	EMAIL      = "j@j.com"
	LOGIN      = "slater100"
	PASSWORD   = "password"
)

var hashedPassword = helpers.IncryptPassword(PASSWORD)
var user = model.User{FirstName: FIRST_NAME, LastName: LAST_NAME, Email: EMAIL, Login: LOGIN, Password: hashedPassword}

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

	want := model.User{Id: 1, FirstName: FIRST_NAME, LastName: LAST_NAME, Email: EMAIL, Login: LOGIN, Password: hashedPassword}
	got, _ := testPostgres.ds.CreateUser(&user)

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("got user: %v but want user: %v", got, &want)
	}
}

func TestCreateUser_idIncrements(t *testing.T) {
	clearTable()

	secondEmail, secondLogin := "second@email.com", "secondLogin"
	secondUser, _ := model.NewUser(FIRST_NAME, LAST_NAME, secondEmail, secondLogin, PASSWORD)

	gotFirst, _ := testPostgres.ds.CreateUser(&user)
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

	testPostgres.ds.CreateUser(&user)

	newUser, _ := model.NewUser("newName", "newName", EMAIL, "newLogin", "newPassword")
	_, err := testPostgres.ds.CreateUser(newUser)

	if err == nil {
		t.Errorf("expected error %v", err)
	}

}

func TestCreateUser_loginNotUnique(t *testing.T) {
	clearTable()

	testPostgres.ds.CreateUser(&user)

	newUser, _ := model.NewUser("newName", "newName", "new@email.com", LOGIN, "newPassword")
	_, err := testPostgres.ds.CreateUser(newUser)

	if err == nil {
		t.Errorf("expected error %v", err)
	}
}

func TestGetUser_success(t *testing.T) {
	clearTable()
	insertUser(user)

	got, _ := testPostgres.ds.GetUser(EMAIL, hashedPassword)
	want := model.User{Id: 1, FirstName: FIRST_NAME, LastName: LAST_NAME, Email: EMAIL}

	if !reflect.DeepEqual(got, &want) {
		t.Errorf("got user: %v but want user: %v", got, &want)
	}
}

func TestGetUser_unknownEmail(t *testing.T) {
	clearTable()
	insertUser(user)

	_, err := testPostgres.ds.GetUser("unknown@email.com", hashedPassword)
	want := fmt.Errorf("unauthorised")

	if err.Error() != want.Error() {
		t.Errorf("got: %v but want error: %v", err.Error(), want.Error())

	}
}

func TestGetUser_unknownPassword(t *testing.T) {
	clearTable()
	insertUser(user)

	_, err := testPostgres.ds.GetUser(EMAIL, "unknownPassword")
	want := fmt.Errorf("unauthorised")

	if err.Error() != want.Error() {
		t.Errorf("got: %v but want error: %v", err.Error(), want.Error())

	}
}

func TestGetUser_unhashedPassword(t *testing.T) {
	clearTable()
	insertUser(user)

	_, err := testPostgres.ds.GetUser(EMAIL, PASSWORD)
	want := fmt.Errorf("unauthorised")

	if err.Error() != want.Error() {
		t.Errorf("got: %v but want error: %v", err.Error(), want.Error())

	}
}

func insertUser(user model.User) {
	err := testPostgres.ds.DB.QueryRow(`INSERT INTO "user" (first_name, last_name, login, email, password_hash) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		user.FirstName, user.LastName, user.Login, user.Email, user.Password).Scan(&user.Id)

	if err != nil {
		log.Panic(err)
	}
}

func clearTable() {
	testPostgres.ds.DB.Exec(`DELETE FROM "user"`)
	testPostgres.ds.DB.Exec(`ALTER SEQUENCE user_id_seq RESTART WITH 1`)
}
