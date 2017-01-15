package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tylerb/graceful"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status"
)

type localEnv struct {
	DB           *sqlx.DB
	IslatestLink *status.SQLIntImpl
}

func main() {

	dbName := "playwithsql"

	dbType := flag.String("db", "mysql", "type of db to bench: mysql, cockroachdb, postgres")
	host := flag.String("host", "localhost", "host IP")
	flag.Parse()

	db, err := dbhandler.Get(*dbType, *host, dbName)
	if err != nil {
		log.Fatalf("%v", err)
	}

	islatestSQLLink := status.GetSQLIntImpl(*dbType)
	err = islatestSQLLink.MigrateUp(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	env := &localEnv{
		DB:           db,
		IslatestLink: islatestSQLLink,
	}

	router := globalMux(env)

	drainInterval, err := time.ParseDuration("2s")
	if err != nil {
		log.Fatalf("%v", err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: ":8083", Handler: router},
	}

	log.Println("Running HTTP server on :8083")
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatalf("%v", err)
	}
}
