package main

import (
	"net/http"
	"fmt"
	"github.com/jack-slater/go-login/app/datastore"
	"log"
	"os"
	"github.com/jack-slater/go-login/app/handlers"
)

func main()  {

	postgresDb := createPostgresDb()
	defer postgresDb.Close()

	createRoutes(postgresDb)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func createRoutes(postgresDb *datastore.PostgresDatastore) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Go-Login")
	})
	http.Handle("/create", handlers.CreateHandler(postgresDb))
	http.Handle("/login", handlers.LoginHandler(postgresDb))
}

func createPostgresDb() *datastore.PostgresDatastore {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		"db", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	postgresDb, err := datastore.NewPostgresDataStore(connectionString)

	if err != nil {
		log.Print(err)
	}
	return postgresDb
}

