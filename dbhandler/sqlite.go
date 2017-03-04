package dbhandler

import (
	"database/sql"
	"os"

	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB is used to have to interface valid
type SQLiteDB struct{}

// NewDBHandler connect to db and return the connection
func (sqliteConf *SQLiteDB) NewDBHandler() (*sql.DB, error) {

	if _, err := os.Stat("./test.db"); os.IsExist(err) {
		errRem := os.Remove("./test.db")
		if errRem != nil {
			return nil, errRem
		}
	}

	return sql.Open("sqlite3", "./test.db")
}
