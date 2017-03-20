package main

import (
	"flag"
	"fmt"
	"log"
	"sort"

	"time"

	"math"

	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status"
)

// "github.com/montanaflynn/stats"
// add median and variance time

// Results to be returned
type Results struct {
	DBType       string
	MaxConns     int
	BenchResults []BenchResult
}

// Bench data
type BenchResult struct {
	Action     string
	Loops      int
	PauseTime  time.Duration
	Errors     int
	Median     time.Duration
	StandDev   time.Duration
	Throughput int
}

func main() {

	// Flags
	dbName := "playwithsql"
	dbType := flag.String("db", "mysql", "type of db to bench: mysql, cockroachdb, postgres")
	dbHost := flag.String("host", "127.0.0.1", "host IP")
	loops := flag.Int("loops", 100, "number of loops")
	pauseTimeStr := flag.String("pausetime", "5ms", "waiting time in ms between each run")
	maxConns := flag.Int("maxconns", 100, "number of max connections")
	flag.Parse()

	pauseTime, err := time.ParseDuration(*pauseTimeStr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	db, err := dbhandler.Get(*dbType, *dbHost, dbName)
	if err != nil {
		log.Fatalf("%s - %s - %s, \n%v", *dbType, *dbHost, dbName, err)
	}

	// Connection
	islatestSQLLink := status.GetSQLIntImpl(*dbType)
	err = islatestSQLLink.MigrateDown(db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = islatestSQLLink.MigrateUp(db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Number of max connections
	// TODO set the param in the db config
	db.SetMaxOpenConns(*maxConns)

	var results = Results{
		DBType:   *dbType,
		MaxConns: *maxConns,
	}

	// Create
	createResults, testEntityoneIDs, err := BenchmarkCreate(*loops, pauseTime, db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, createResults)

	// Update
	updateResults, err := BenchmarkUpdateStatus(*loops, pauseTime, db, islatestSQLLink, testEntityoneIDs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, updateResults)

	// Select by status
	selectByStatusResults, err := BenchmarkSelectEntityoneByStatus(*loops, db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, selectByStatusResults)

	// Select by PK
	selectByPKResults, err := BenchmarkSelectEntityoneOneByPK(*loops, db, islatestSQLLink, testEntityoneIDs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, selectByPKResults)

	jsonResults, err := json.Marshal(results)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("%s", jsonResults)
}

// BenchmarkCreate will loop a loops number of time and give the resulting time taken
func BenchmarkCreate(
	loops int, pauseTime time.Duration, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl,
) (
	results BenchResult,
	testEntityoneIDs []int64,
	err error,
) {
	type latID struct {
		dur time.Duration
		id  int64
	}
	latenciesIDsC := make(chan latID)
	entityIDsC := make(chan int64)
	errorC := make(chan error)

	defer close(latenciesIDsC)
	defer close(entityIDsC)
	defer close(errorC)

	var latencies []time.Duration
	var errCount int

	before := time.Now()

	for i := 0; i < loops; i++ {
		time.Sleep(pauseTime)
		go func() {
			var e status.Entityone
			before := time.Now()
			errCr := e.Create(dbConn, benchSQLLink)
			if errCr != nil {
				errorC <- errCr
			} else {
				latenciesIDsC <- latID{dur: time.Now().Sub(before), id: e.ID}
			}
		}()
	}

	for j := 0; j < loops; j++ {
		select {
		case latencyID := <-latenciesIDsC:
			latencies = append(latencies, latencyID.dur)
			testEntityoneIDs = append(testEntityoneIDs, latencyID.id)
		case errCr := <-errorC:
			log.Printf("%v", errCr)
			errCount++
		}
	}

	timeTaken := time.Now().Sub(before)

	return BenchResult{
			Action:     "create",
			Loops:      loops,
			PauseTime:  pauseTime,
			Errors:     errCount,
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(timeTaken.Seconds() / float64(loops)),
		},
		testEntityoneIDs,
		nil
}

// BenchmarkUpdateStatus benchmark for status updates (include deletes)
func BenchmarkUpdateStatus(
	loops int, pauseTime time.Duration, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl, testEntityoneIDs []int64,
) (
	results BenchResult,
	err error,
) {
	if len(testEntityoneIDs) == 0 {
		return results, fmt.Errorf("BenchmarkUpdateStatus: no entity created, nothing to update")
	}

	latenciesC := make(chan time.Duration)
	errorC := make(chan error)

	defer close(latenciesC)
	defer close(errorC)

	var latencies []time.Duration
	var errCount int

	before := time.Now()

	for i := 0; i < loops; i++ {
		time.Sleep(pauseTime)
		go func() {
			var e status.Entityone
			e.ID = testEntityoneIDs[i%len(testEntityoneIDs)]

			before := time.Now()
			errU := e.UpdateStatus(dbConn, benchSQLLink, status.ActionCancel, status.StatusCancelled)
			if errU != nil {
				errorC <- errU
			} else {
				latenciesC <- time.Now().Sub(before)
			}
		}()
	}

	for j := 0; j < loops; j++ {
		select {
		case latency := <-latenciesC:
			latencies = append(latencies, latency)
		case errU := <-errorC:
			log.Printf("%v", errU)
			errCount++
		}
	}

	timeTaken := time.Now().Sub(before)

	return BenchResult{
			Action:     "create",
			Loops:      loops,
			PauseTime:  pauseTime,
			Errors:     errCount,
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(timeTaken.Seconds() / float64(loops)),
		},
		nil

}

// BenchmarkSelectEntityoneByStatus benchmark with select by status
func BenchmarkSelectEntityoneByStatus(
	loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl,
) (
	results BenchResult,
	err error,
) {
	latenciesC := make(chan time.Duration)
	errorC := make(chan error)

	defer close(latenciesC)
	defer close(errorC)

	var latencies []time.Duration
	var errCount int

	before := time.Now()

	go func() {
		for i := 0; i < loops; i++ {
			_, errSel := status.SelectEntityoneByStatus(dbConn, benchSQLLink, status.StatusCancelled)
			if errSel != nil {
				errorC <- errSel
			} else {
				latenciesC <- time.Now().Sub(before)
			}
		}
	}()

	for j := 0; j < loops; j++ {
		select {
		case latency := <-latenciesC:
			latencies = append(latencies, latency)
		case errSel := <-errorC:
			log.Printf("%v", errSel)
			errCount++
		}
	}

	timeTaken := time.Now().Sub(before)

	return BenchResult{
			Action:     "create",
			Loops:      loops,
			PauseTime:  0,
			Errors:     errCount,
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(timeTaken.Seconds() / float64(loops)),
		},
		nil
}

// BenchmarkSelectEntityoneOneByPK benchmark with select by primary key
func BenchmarkSelectEntityoneOneByPK(
	loops int, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl, testEntityoneIDs []int64,
) (
	results BenchResult,
	err error,
) {
	latenciesC := make(chan time.Duration)
	errorC := make(chan error)

	defer close(latenciesC)
	defer close(errorC)

	var latencies []time.Duration
	var errCount int

	before := time.Now()

	go func() {
		for i := 0; i < loops; i++ {
			_, errSel := status.SelectEntityoneOneByPK(dbConn, benchSQLLink, testEntityoneIDs[i%len(testEntityoneIDs)])
			if errSel != nil {
				errorC <- errSel
			} else {
				latenciesC <- time.Now().Sub(before)
			}

		}
	}()

	for j := 0; j < loops; j++ {
		select {
		case latency := <-latenciesC:
			latencies = append(latencies, latency)
		case errSel := <-errorC:
			log.Printf("%v", errSel)
			errCount++
		}
	}

	timeTaken := time.Now().Sub(before)

	return BenchResult{
			Action:     "create",
			Loops:      loops,
			PauseTime:  0,
			Errors:     errCount,
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(timeTaken.Seconds() / float64(loops)),
		},
		nil
}

// returns the median duration of a list
func getMedian(latencies []time.Duration) time.Duration {
	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	if len(latencies) == 0 {
		return 0
	}
	if len(latencies) == 1 {
		return latencies[0]
	}
	if len(latencies)%2 == 0 {
		return latencies[(len(latencies)/2-1)] + latencies[(len(latencies)/2+1)]
	}
	return latencies[len(latencies)/2]
}

// returns the standard deviation of the list
func getStandardDeviation(latencies []time.Duration) time.Duration {

	if len(latencies) == 0 {
		return 0
	}

	// Sum the square of the mean subtracted from each number
	mean := getMean(latencies)

	var variance time.Duration

	for _, latency := range latencies {
		variance += (latency - mean) * (latency - mean)
	}

	return time.Duration(math.Pow(float64(variance.Nanoseconds()/int64(len(latencies))), 0.5))
}

func getMean(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}

	return time.Duration(total.Nanoseconds() / int64(len(latencies)))
}
