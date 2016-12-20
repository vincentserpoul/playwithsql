package dbhandler

// import (
// 	"github.com/serenize/snaker"

// 	"github.com/jmoiron/sqlx"

// 	// to connect to Oracle
// 	_ "gopkg.in/rana/ora.v3"
// )

// // OracleDB is a conf for the mysql database
// type OracleDB struct {
// 	Host     string
// 	Port     string
// 	User     string
// 	Password string
// 	Sid      string
// }

// // NewDBHandler connect to db and return the connection
// func (OracleConf OracleDB) NewDBHandler() (*sqlx.DB, error) {

// 	dsn := OracleConf.User + "/"
// 	OracleConf.Password + "@"
// 	OracleConf.Host + ":"
// 	OracleConf.Port + "/"
// 	OracleConf.Sid

// 	db := sqlx.MustConnect("ora", dsn)

// 	db.MapperFunc(snaker.CamelToSnake)

// 	return db, nil
// }
