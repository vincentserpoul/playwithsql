package main

import (
	"flag"
	"log"
	"strconv"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status"
)

func main() {

	dbName := "playwithsql"
	dbType := flag.String("db", "mysql", "type of db to bench: mysql, cockroachdb, postgres")
	dbHost := flag.String("host", "localhost", "host IP")
	loopsStr := flag.String("loops", "10000", "number of loops")
	flag.Parse()

	loops, err := strconv.ParseInt(*loopsStr, 10, 64)
	if err != nil {
		log.Fatalf("%v", err)
	}

	db, err := dbhandler.Get(*dbType, *dbHost, dbName)
	if err != nil {
		log.Fatalf("%v", err)
	}

	islatestSQLLink := status.GetSQLLinkContainer(*dbType)
	err = islatestSQLLink.MigrateDown(db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = islatestSQLLink.MigrateUp(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	createTimeTaken, testEntityoneIDs, err := BenchmarkCreate(int(loops), db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%s;create;%d;%d;%d", *dbType, loops, createTimeTaken.Nanoseconds(), len(testEntityoneIDs))
}

// BenchmarkCreate will loop a loops number of time and give the resulting time taken
func BenchmarkCreate(loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLLinkContainer) (
	timeTaken time.Duration,
	testEntityoneIDs []int64,
	err error,
) {
	before := time.Now()
	for i := 0; i < loops; i++ {
		var e status.Entityone
		errCr := e.Create(dbConn, benchSQLLink)
		if errCr != nil {
			return timeTaken, testEntityoneIDs, errCr
		}
		i++
		// limit the number of entities to store
		if len(testEntityoneIDs) < 500 {
			testEntityoneIDs = append(testEntityoneIDs, e.ID)
		}
	}

	after := time.Now()
	timeTaken = after.Sub(before)

	return timeTaken, testEntityoneIDs, err
}

// func BenchmarkUpdateStatus(b *testing.B) {
// 	var e Entityone
// 	for i := 0; i < b.N; i++ {
// 		e.ID = testEntityoneIDs[b.N%len(testEntityoneIDs)]
// 		_ = e.UpdateStatus(testDBConn, testSQLLink, ActionCancel, StatusCancelled)
// 	}
// }

// func BenchmarkSelectEntityoneByStatus(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, err := SelectEntityoneByStatus(testDBConn, testSQLLink, StatusCreated)
// 		if err != nil {
// 			b.Errorf("Select entityone by status: %v", err)
// 			return
// 		}
// 	}
// }

// func BenchmarkSelectEntityoneOneByPK(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		_, err := SelectEntityoneOneByPK(testDBConn, testSQLLink, testEntityoneIDs[b.N%len(testEntityoneIDs)])
// 		if err != nil {
// 			b.Errorf("Select entityone by status: %v", err)
// 			return
// 		}
// 	}
// }
