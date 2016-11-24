package dbhandler

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/serenize/snaker"
)

// MySQLDB is a conf for the mysql database
type MySQLDB struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	SSL      SSL
}

// SSL is a conf used to connect to SSL if needed
type SSL struct {
	CertPath   string
	KeyPath    string
	CAPath     string
	ServerName string
}

// NewDBHandler connect to db and return the connection
func (mysqlConf *MySQLDB) NewDBHandler() (*sqlx.DB, error) {
	dsn := mysqlConf.User + ":" +
		mysqlConf.Password + "@" +
		mysqlConf.Protocol + "(" +
		mysqlConf.Host + ":" +
		mysqlConf.Port + ")/" +
		mysqlConf.Dbname + "?parseTime=true"

	if mysqlConf.SSL.CAPath != "" &&
		mysqlConf.SSL.CertPath != "" &&
		mysqlConf.SSL.KeyPath != "" {

		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(mysqlConf.SSL.CAPath)
		if err != nil {
			return nil, fmt.Errorf("infrastructure NewDBHandler: %v", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, fmt.Errorf("infrastructure NewDBHandler: Failed to append PEM")
		}
		clientCert := make([]tls.Certificate, 0, 1)
		certs, err := tls.LoadX509KeyPair(mysqlConf.SSL.CertPath, mysqlConf.SSL.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("infrastructure NewDBHandler: %v", err)
		}
		clientCert = append(clientCert, certs)
		err = mysql.RegisterTLSConfig("mysqltls", &tls.Config{
			RootCAs:      rootCertPool,
			Certificates: clientCert,
			ServerName:   mysqlConf.SSL.ServerName,
		})
		if err != nil {
			return nil, fmt.Errorf("infrastructure NewDBHandler: %v", err)
		}
		dsn += "&tls=mysqltls"
	}

	db := sqlx.MustConnect("mysql", dsn)

	db.MapperFunc(snaker.CamelToSnake)

	// we don't want to close the connection

	//  if err != nil {
	//      log.Printf("infrastructure db Close(): %v", err)
	//  }
	// }()

	return db, nil
}
