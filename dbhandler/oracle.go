package dbhandler

import

// to connect to Oracle
(
	"database/sql"

	_ "gopkg.in/rana/ora.v4"
)

// OracleDB is a conf for the mysql database
type OracleDB struct {
	Host     string
	Port     string
	User     string
	Password string
	Sid      string
}

// NewDBHandler connect to db and return the connection
func (OracleConf OracleDB) NewDBHandler() (*sql.DB, error) {

	dsn := OracleConf.User + "/" +
		OracleConf.Password + "@" +
		OracleConf.Host + ":" +
		OracleConf.Port + "/" +
		OracleConf.Sid

	return sql.Open("ora", dsn)
}
