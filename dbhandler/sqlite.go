package dbhandler

import (
	"os"

	"github.com/jmoiron/sqlx"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB is used to have to interface valid
type SQLiteDB struct{}

// NewDBHandler connect to db and return the connection
func (sqliteConf *SQLiteDB) NewDBHandler() (*sqlx.DB, error) {

	if _, err := os.Stat("./test.db"); os.IsExist(err) {
		errRem := os.Remove("./test.db")
		if errRem != nil {
			return nil, errRem
		}
	}

	db := sqlx.MustConnect("sqlite3", "./test.db")

	return db, nil
}
