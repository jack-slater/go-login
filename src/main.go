package main

import (
	"net/http"
	"fmt"
	"github.com/jack-slater/go-login/src/datastore"
	"log"
	"github.com/jack-slater/go-login/src/handlers"
)

func main()  {


	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", "postgres-dev","password", "dev")
	postgresDb, err := datastore.NewPostgresDataStore(connectionString)

	if err != nil {
		log.Fatal(err)
	}

	defer postgresDb.Close()

	http.Handle("/", handlers.UserHandler(postgresDb))

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

