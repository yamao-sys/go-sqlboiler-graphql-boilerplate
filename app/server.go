package main

import (
	"app/db"
	"app/lib"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// NOTE: DB接続
	dbCon := db.Init()
	isDebug, _ := strconv.ParseBool(os.Getenv("DB_DEBUG_MODE"))
	boil.DebugMode = isDebug

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", lib.GetGraphQLHttpHandler(dbCon))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
