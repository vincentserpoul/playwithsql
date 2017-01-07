package main

import (
	"flag"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status"
)

type localEnv struct {
	DB           *sqlx.DB
	IslatestLink *status.SQLLinkContainer
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

	islatestSQLLink := status.GetSQLLinkContainer(*dbType)
	err = islatestSQLLink.MigrateUp(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	env := &localEnv{
		DB:           db,
		IslatestLink: islatestSQLLink,
	}

	log.Printf("Running bench! db: %v", env.DB.Stats())

}
