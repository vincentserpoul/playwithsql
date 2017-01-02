package dbhandler

import (
	"github.com/jmoiron/sqlx"

	// to connect to mssql
	_ "github.com/denisenkom/go-mssqldb"
)

// MSSQLDB is a conf for the mysql database
type MSSQLDB struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

// NewDBHandler connect to db and return the connection
func (MSSQLConf MSSQLDB) NewDBHandler() (*sqlx.DB, error) {

	dsn := "server=" + MSSQLConf.Host +
		";port=" + MSSQLConf.Port +
		";user id=" + MSSQLConf.User +
		";password=" + MSSQLConf.Password +
		";database=" + MSSQLConf.Dbname

	db := sqlx.MustConnect("mssql", dsn)

	return db, nil
}
