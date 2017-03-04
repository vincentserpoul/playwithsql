package dbhandler

import (
	"database/sql"

	// to connect to cockroachdb
	_ "github.com/lib/pq"
)

// CockroachDB is a conf for the mysql database
type CockroachDB struct {
	Host   string
	Port   string
	User   string
	Dbname string
	SSL    SSL
}

// NewDBHandler connect to db and return the connection
func (CockroachConf CockroachDB) NewDBHandler() (*sql.DB, error) {

	dsn := "postgresql://" +
		CockroachConf.User + "@" +
		CockroachConf.Host + ":" +
		CockroachConf.Port + "/" +
		CockroachConf.Dbname + "?sslmode=disable"

	return sql.Open("postgres", dsn)
}
