package dbhandler

import (
	"os"

	"github.com/jmoiron/sqlx"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/serenize/snaker"
)

// SQLiteDB is used to have to interface valid
type SQLiteDB struct{}

// NewDBHandler connect to db and return the connection
func (sqliteConf *SQLiteDB) NewDBHandler() (*sqlx.DB, error) {
	_ = os.Remove("./test.db")

	db := sqlx.MustConnect("sqlite3", "./test.db")

	db.MapperFunc(snaker.CamelToSnake)

	return db, nil
}
