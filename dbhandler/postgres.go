package dbhandler

import

// to connect to cockroachdb
(
	"database/sql"

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
func (PostgresConf PostgresDB) NewDBHandler() (*sql.DB, error) {

	dsn := "postgres://" +
		PostgresConf.User + ":" +
		PostgresConf.Password + "@" +
		PostgresConf.Host + ":" +
		PostgresConf.Port + "/" +
		PostgresConf.Dbname + "?sslmode=disable"

	return sql.Open("postgres", dsn)
}
