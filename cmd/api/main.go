package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sidharth-rashwana/book/internal/database"
	environment "github.com/sidharth-rashwana/book/internal/environment"
	routes "github.com/sidharth-rashwana/book/internal/route"
)

func main() {
	uri, dbName, port, host := environment.InitalizeEnv()

	db, client, err := database.InitializeDB(uri, dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting MongoDB client: %v", err)
		}
	}()

	r := routes.InitalizeRoutes(db)

	log.Println("Server is starting and listening on " + host + " on port: " + port)

	log.Fatal(http.ListenAndServe(host+":"+port, r))
}
