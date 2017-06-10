package dbhandler

import (
	"github.com/jmoiron/sqlx"
	// to connect to cockroachdb
	_ "github.com/lib/pq"
)

// CockroachDB is a conf for the mysql database
type CockroachDB struct {
	Host   string
	Port   string
	User   string
	Dbname string
}

// NewDBHandler connect to db and return the connection
func (CockroachConf CockroachDB) NewDBHandler() (*sqlx.DB, error) {

	dsn := "postgresql://" +
		CockroachConf.User + "@" +
		CockroachConf.Host + ":" +
		CockroachConf.Port + "/" +
		CockroachConf.Dbname + "?sslmode=disable"

	db := sqlx.MustConnect("postgres", dsn)

	return db, nil
}
