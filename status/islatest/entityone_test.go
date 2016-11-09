package islatest

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql"
	"github.com/vincentserpoul/playwithsql/status/islatest/cockroachdb"
	"github.com/vincentserpoul/playwithsql/status/islatest/mysql"
	"github.com/vincentserpoul/playwithsql/status/islatest/postgres"
)

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var e Entityone
		_ = e.Create(db, sqlLink)
	}
}

func TestCreate(t *testing.T) {
	var e Entityone
	err := e.Create(db, sqlLink)
	if err != nil {
		t.Errorf("create entityone: %v", err)
	}
}

func TestUpdateStatus(t *testing.T) {
	var e Entityone

	err := e.Create(db, sqlLink)
	if err != nil {
		t.Errorf("UpdateStatus entityone: %v", err)
	}
	err = e.UpdateStatus(db, sqlLink, ActionCancel, StatusCancelled)
	if err != nil {
		t.Errorf("UpdateStatus entityone: %v", err)
	}

	if e.ActionID != ActionCancel && e.StatusID != StatusCancelled {
		t.Errorf("UpdateStatus entityone: status and action not updated")
	}
}

func BenchmarkUpdateStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = entityToUpdate.UpdateStatus(db, sqlLink, ActionCancel, StatusCancelled)
	}
}

func TestSelect(t *testing.T) {
	var e Entityone

	err := e.Create(db, sqlLink)
	if err != nil {
		t.Errorf("Select entityone: %v", err)
	}

	entityOnes, err := SelectEntityone(db, sqlLink)
	if err != nil {
		t.Errorf("Select entityone: %v", err)
	}

	if len(entityOnes) == 0 {
		t.Errorf("Select entityone didn't retrieve any")
	}
}

func BenchmarkSelect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SelectEntityone(db, sqlLink)
	}
}

var db *sqlx.DB
var sqlLink SQLLink
var entityToUpdate Entityone

func TestMain(m *testing.M) {

	var err error
	var conf playwithsql.ConfType
	dbName := "entityone_test"

	dbType := flag.String("db", "mysql", "type of db to bench: mysql, cockroachdb")
	flag.Parse()

	switch *dbType {
	case "mysql":
		conf = &playwithsql.MySQLDB{
			Protocol: "tcp",
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "test",
			Dbname:   dbName,
			SSL: playwithsql.SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
		sqlLink = &mysql.Link{}
	case "postgres":
		conf = &playwithsql.PostgresDB{
			Host:     "localhost",
			Port:     "5432",
			User:     "root",
			Password: "test",
			Dbname:   dbName,
			SSL: playwithsql.SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
		sqlLink = &postgres.Link{}
	case "cockroachdb":
		conf = &playwithsql.CockroachDB{
			Host:   "localhost",
			Port:   "26257",
			User:   "root",
			Dbname: dbName,
			SSL: playwithsql.SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
		sqlLink = &cockroachdb.Link{}
	}

	db, err = conf.NewDBHandler()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = sqlLink.MigrateDown(db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = sqlLink.MigrateUp(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// We need at least one entity for the update
	err = entityToUpdate.Create(db, sqlLink)
	if err != nil {
		log.Fatalf("%v", err)
	}

	retCode := m.Run()

	err = sqlLink.MigrateDown(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	os.Exit(retCode)
}
