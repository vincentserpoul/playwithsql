package dbhandler

import (
	"github.com/jmoiron/sqlx"

	// to connect to mysql \
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDB is a conf for the mysql database
type MySQLDB struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

// NewDBHandler connect to db and return the connection
func (mysqlConf *MySQLDB) NewDBHandler() (*sqlx.DB, error) {
	dsn := mysqlConf.User + ":" +
		mysqlConf.Password + "@" +
		mysqlConf.Protocol + "(" +
		mysqlConf.Host + ":" +
		mysqlConf.Port + ")/" +
		mysqlConf.Dbname + "?parseTime=true&multiStatements=true"

	db := sqlx.MustConnect("mysql", dsn)

	return db, nil
}
