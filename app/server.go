package main

import (
	"app/db"
	"app/lib"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// NOTE: DB接続
	dbCon := db.Init()

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", lib.GetGraphQLHttpHandler(dbCon))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
