package main

import (
	"flag"
	"fmt"
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

	islatestSQLLink := status.GetSQLIntImpl(*dbType)
	err = islatestSQLLink.MigrateDown(db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = islatestSQLLink.MigrateUp(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create
	createTimeTaken, testEntityoneIDs, err := BenchmarkCreate(int(loops), db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%s;create;%d;%d\n", *dbType, loops, int64(loops/int64(createTimeTaken.Seconds())))

	// Update
	updateTimeTaken, err := BenchmarkUpdateStatus(int(loops), db, islatestSQLLink, testEntityoneIDs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%s;update;%d;%d\n", *dbType, loops, int64(loops/int64(updateTimeTaken.Seconds())))

	// Select by status
	selectByStatusTimeTaken, err := BenchmarkSelectEntityoneByStatus(int(loops), db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%s;selectByStatus;%d;%d\n", *dbType, loops, int64(loops/int64(selectByStatusTimeTaken.Seconds())))

	// Select by PK
	selectByPKTimeTaken, err := BenchmarkSelectEntityoneOneByPK(int(loops), db, islatestSQLLink, testEntityoneIDs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%s;selectByPK;%d;%v\n", *dbType, loops, int64(loops/int64(selectByPKTimeTaken.Seconds())))
}

// BenchmarkCreate will loop a loops number of time and give the resulting time taken
func BenchmarkCreate(loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl) (
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

// BenchmarkUpdateStatus benchmark for status updates (include deletes)
func BenchmarkUpdateStatus(loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl, testEntityoneIDs []int64) (
	timeTaken time.Duration,
	err error,
) {
	var e status.Entityone
	before := time.Now()

	for i := 0; i < loops; i++ {
		e.ID = testEntityoneIDs[i%len(testEntityoneIDs)]
		errU := e.UpdateStatus(dbConn, benchSQLLink, status.ActionCancel, status.StatusCancelled)
		if errU != nil {
			return timeTaken, errU
		}
	}

	after := time.Now()
	timeTaken = after.Sub(before)

	return timeTaken, err
}

// BenchmarkSelectEntityoneByStatus benchmark with select by status
func BenchmarkSelectEntityoneByStatus(loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl) (
	timeTaken time.Duration,
	err error,
) {
	before := time.Now()
	for i := 0; i < loops; i++ {
		_, errSel := status.SelectEntityoneByStatus(dbConn, benchSQLLink, status.StatusCancelled)
		if errSel != nil {
			return timeTaken, errSel
		}
	}

	after := time.Now()
	timeTaken = after.Sub(before)

	return timeTaken, nil
}

// BenchmarkSelectEntityoneOneByPK benchmark with select by primary key
func BenchmarkSelectEntityoneOneByPK(loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl, testEntityoneIDs []int64) (
	timeTaken time.Duration,
	err error,
) {
	before := time.Now()
	for i := 0; i < loops; i++ {
		_, errSel := status.SelectEntityoneOneByPK(dbConn, benchSQLLink, testEntityoneIDs[i%len(testEntityoneIDs)])
		if err != nil {
			return timeTaken, errSel
		}
	}
	after := time.Now()
	timeTaken = after.Sub(before)

	return timeTaken, nil
}
