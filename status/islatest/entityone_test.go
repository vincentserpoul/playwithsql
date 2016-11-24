package islatest

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status/islatest/cockroachdb"
	"github.com/vincentserpoul/playwithsql/status/islatest/mysql"
	"github.com/vincentserpoul/playwithsql/status/islatest/postgres"
	"github.com/vincentserpoul/playwithsql/status/islatest/sqlite"
)

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var e Entityone
		_ = e.Create(db, sqlLink)

		// limit the number of tests
		if len(testEntityoneIDs) < 500 {
			testEntityoneIDs = append(testEntityoneIDs, e.ID)
		}
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
	var e Entityone
	for i := 0; i < b.N; i++ {
		e.ID = testEntityoneIDs[b.N%len(testEntityoneIDs)]
		_ = e.UpdateStatus(db, sqlLink, ActionCancel, StatusCancelled)
	}
}

func TestSelectEntityoneOneByStatus(t *testing.T) {
	var e Entityone

	err := e.Create(db, sqlLink)
	if err != nil {
		t.Errorf("Select entityone by status: %v", err)
	}

	_, err = SelectEntityoneOneByStatus(db, sqlLink, StatusCreated)
	if err != nil {
		t.Errorf("Select entityone by status: %v", err)
	}
}

func BenchmarkSelectEntityoneOneByStatus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := SelectEntityoneOneByStatus(db, sqlLink, StatusCreated)
		if err != nil {
			b.Errorf("Select entityone by status: %v", err)
			return
		}
	}
}

func TestSelectEntityoneOneByPK(t *testing.T) {
	var e Entityone

	err := e.Create(db, sqlLink)
	if err != nil {
		t.Errorf("Select entityone by pk: %v", err)
		return
	}

	entityOne, err := SelectEntityoneOneByPK(db, sqlLink, e.ID)
	if err != nil {
		t.Errorf("Select entityone by pk: %v", err)
		return
	}

	if entityOne.ID != e.ID {
		t.Errorf("Select entityone by pk retrieved entity %d instead of %d", entityOne.ID, e.ID)
		return
	}
	var emptyTime time.Time
	if entityOne.TimeCreated == emptyTime {
		t.Errorf("Select entityone by pk retrieved but entity time created not correctly retrieved: %v", entityOne)
		return
	}
	if entityOne.Status.TimeCreated == emptyTime {
		t.Errorf("Select entityone by pk retrieved but entity status time created not correctly retrieved: %v", entityOne)
		return
	}
	if entityOne.Status.ActionID == 0 {
		t.Errorf("Select entityone by pk retrieved but entity status actionid created not correctly retrieved: %v", entityOne)
		return
	}
	if entityOne.Status.StatusID == 0 {
		t.Errorf("Select entityone by pk retrieved but entity status statusid created not correctly retrieved: %v", entityOne)
		return
	}
}

func BenchmarkSelectEntityoneOneByPK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := SelectEntityoneOneByPK(db, sqlLink, testEntityoneIDs[b.N%len(testEntityoneIDs)])
		if err != nil {
			b.Errorf("Select entityone by status: %v", err)
			return
		}
	}
}

var db *sqlx.DB
var sqlLink SQLLink
var testEntityoneIDs []int64

func TestMain(m *testing.M) {

	var err error
	var conf dbhandler.ConfType
	dbName := "entityone_test"

	dbType := flag.String("db", "mysql", "type of db to bench: mysql, cockroachdb, postgres")
	host := flag.String("host", "localhost", "host IP")
	flag.Parse()

	switch *dbType {
	case "mysql":
		conf = &dbhandler.MySQLDB{
			Protocol: "tcp",
			Host:     *host,
			Port:     "3306",
			User:     "root",
			Password: "test",
			Dbname:   dbName,
			SSL: dbhandler.SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
		sqlLink = &mysql.Link{}
	case "sqlite":
		conf = &dbhandler.SQLiteDB{}
		sqlLink = &sqlite.Link{}
	case "postgres":
		conf = &dbhandler.PostgresDB{
			Host:     *host,
			Port:     "5432",
			User:     "root",
			Password: "test",
			Dbname:   dbName,
			SSL: dbhandler.SSL{
				CertPath:   "",
				KeyPath:    "",
				CAPath:     "",
				ServerName: "",
			},
		}
		sqlLink = &postgres.Link{}
	case "cockroachdb":
		conf = &dbhandler.CockroachDB{
			Host:   "localhost",
			Port:   "26257",
			User:   "root",
			Dbname: dbName,
			SSL: dbhandler.SSL{
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
	defer func() {
		errClose := db.Close()
		if errClose != nil {
			log.Fatalf("%v", errClose)
		}
	}()

	err = sqlLink.MigrateDown(db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = sqlLink.MigrateUp(db)
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
