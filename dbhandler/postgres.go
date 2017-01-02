package dbhandler

import (
	"github.com/jmoiron/sqlx"
	// to connect to cockroachdb
	_ "github.com/lib/pq"
)

// PostgresDB is a conf for the mysql database
type PostgresDB struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSL      SSL
}

// NewDBHandler connect to db and return the connection
func (PostgresConf PostgresDB) NewDBHandler() (*sqlx.DB, error) {

	dsn := "postgres://" +
		PostgresConf.User + ":" +
		PostgresConf.Password + "@" +
		PostgresConf.Host + ":" +
		PostgresConf.Port + "/" +
		PostgresConf.Dbname + "?sslmode=disable"

	db := sqlx.MustConnect("postgres", dsn)

	return db, nil
}
