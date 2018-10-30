package datastore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/jack-slater/go-login/app/model"
	"log"
)

type PostgresDatastore struct {
	*sql.DB
}

func NewPostgresDataStore(connection string) (*PostgresDatastore, error) {

	connectedDb, err := sql.Open("postgres", connection)
	if err != nil {
		return nil, err
	}
	datastore := &PostgresDatastore{connectedDb}
	datastore.createUserTable()
	return datastore, nil
}

func (p *PostgresDatastore) GetUser(login, password string) error {
	return nil
}

func (p *PostgresDatastore) CreateUser(user *model.User) (*model.User, error) {

	err := p.DB.QueryRow(`INSERT INTO "user" (first_name, last_name, login, email, password_hash) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		user.FirstName, user.LastName, user.Login, user.Email, user.Password).Scan(&user.Id)

	if err != nil {
		log.Printf("Unable to save user: %+v due to error: %v", user, err)
		return nil, err
	}

	return user, nil
}

func (p *PostgresDatastore) Close() {
	p.DB.Close()
}

func (p *PostgresDatastore) createUserTable() {
	const qry = `CREATE TABLE IF NOT EXISTS "user" ( 
id SERIAL PRIMARY KEY,
	first_name text NOT NULL,
	last_name text NOT NULL,
	email text NOT NULL UNIQUE,
	login text NOT NULL UNIQUE ,
	password_hash text NOT NULL
)`

	if _, err := p.DB.Exec(qry); err != nil {
		log.Fatal(err)
	}
}
