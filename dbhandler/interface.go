package dbhandler

import "database/sql"

// ConfType contains the creation of the handler
type ConfType interface {
	NewDBHandler() (*sql.DB, error)
}

// Get returns the DB connection
func Get(dbType string, host string, dbName string) (db *sql.DB, err error) {

	var conf ConfType

	switch dbType {
	case "mysql", "percona", "mariadb":
		conf = &MySQLDB{
			Protocol: "tcp",
			Host:     host,
			Port:     "3306",
			User:     "root",
			Password: "test",
			Dbname:   dbName,
			SSL: SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
	case "sqlite":
		conf = &SQLiteDB{}
	case "postgres":
		conf = &PostgresDB{
			Host:     host,
			Port:     "5432",
			User:     "root",
			Password: "test",
			Dbname:   dbName,
			SSL: SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
	case "cockroachdb":
		conf = &CockroachDB{
			Host:   host,
			Port:   "26257",
			User:   "root",
			Dbname: dbName,
			SSL: SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
	case "mssql":
		conf = &MSSQLDB{
			Host:     host,
			Port:     "1433",
			User:     "sa",
			Password: "thank5MsSQLforcingMe",
			Dbname:   dbName,
		}
	case "oracle":
		conf = &OracleDB{
			Host:     host,
			Port:     "1521",
			User:     dbName,
			Password: "dev",
			Sid:      "XE",
		}
	}

	return conf.NewDBHandler()
}
