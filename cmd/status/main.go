package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"sync"

	"time"

	"math"

	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status"
)

// Results to be returned
type Results struct {
	DBType       string
	MaxConns     int
	Date         time.Time
	BenchResults []BenchResult
}

// BenchResult data
type BenchResult struct {
	Action     string
	Loops      int
	PauseTime  time.Duration
	Errors     int
	Min        time.Duration
	Max        time.Duration
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
		Date:     time.Now(),
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

	fmt.Printf("%s\n", jsonResults)
}

// BenchmarkCreate will loop a loops number of time and give the resulting time taken
func BenchmarkCreate(
	loops int, pauseTime time.Duration, dbConn *sqlx.DB, benchSQLLink *status.SQLIntImpl,
) (
	results BenchResult,
	testEntityoneIDs []int64,
	err error,
) {
	latenciesC := make(chan time.Duration)
	entityIDsC := make(chan int64)
	errorC := make(chan error)

	defer close(latenciesC)
	defer close(entityIDsC)
	defer close(errorC)

	before := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < loops; i++ {
		time.Sleep(pauseTime)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var e status.Entityone
			beforeLocal := time.Now()
			ok := false
			var errCr error
			retryCount := 0
			for retryCount < 3 && !ok {
				errCr = e.Create(dbConn, benchSQLLink)
				if errCr != nil {
					retryCount++
					time.Sleep(pauseTime * time.Duration(retryCount*10))
				} else {
					ok = true
				}
			}
			if errCr != nil {
				errorC <- errCr
			} else {
				latenciesC <- time.Since(beforeLocal)
				entityIDsC <- e.ID
			}
		}(&wg)
	}

	var latencies []time.Duration
	var errCount int
	go receiveResults(&latenciesC, &errorC, &latencies, &errCount)

	// Receive the entityIDs
	go func() {
		for entityID := range entityIDsC {
			testEntityoneIDs = append(testEntityoneIDs, entityID)
		}
	}()

	wg.Wait()
	timeTaken := time.Since(before)

	return BenchResult{
			Action:     "create",
			Loops:      loops,
			PauseTime:  pauseTime,
			Errors:     errCount,
			Min:        getMin(latencies),
			Max:        getMax(latencies),
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(float64(loops) / timeTaken.Seconds()),
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

	before := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < loops; i++ {
		time.Sleep(pauseTime / time.Duration(2))
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			var e status.Entityone
			e.ID = testEntityoneIDs[i%len(testEntityoneIDs)]
			beforeLocal := time.Now()
			ok := false
			var errU error
			retryCount := 0
			for retryCount < 3 && !ok {
				errU = e.UpdateStatus(dbConn, benchSQLLink, status.ActionCancel, status.StatusCancelled)
				if errU != nil {
					retryCount++
					time.Sleep(pauseTime * time.Duration(retryCount*10))
				} else {
					ok = true
				}
			}
			if errU != nil {
				errorC <- errU
			} else {
				latenciesC <- time.Since(beforeLocal)
			}
		}(&wg)
	}

	var latencies []time.Duration
	var errCount int
	go receiveResults(&latenciesC, &errorC, &latencies, &errCount)

	wg.Wait()
	timeTaken := time.Since(before)

	return BenchResult{
			Action:     "updateStatus",
			Loops:      loops,
			PauseTime:  pauseTime,
			Errors:     errCount,
			Min:        getMin(latencies),
			Max:        getMax(latencies),
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(float64(loops) / timeTaken.Seconds()),
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

	var wg sync.WaitGroup

	before := time.Now()

	for i := 0; i < loops; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			beforeLocal := time.Now()
			_, errSel := status.SelectEntityoneByStatus(dbConn, benchSQLLink, status.StatusCancelled)
			if errSel != nil {
				errorC <- errSel
			} else {
				latenciesC <- time.Since(beforeLocal)
			}
		}(&wg)
	}

	var latencies []time.Duration
	var errCount int
	go receiveResults(&latenciesC, &errorC, &latencies, &errCount)

	wg.Wait()
	timeTaken := time.Since(before)

	return BenchResult{
			Action:     "selectEntityoneByStatus",
			Loops:      loops,
			PauseTime:  0,
			Errors:     errCount,
			Min:        getMin(latencies),
			Max:        getMax(latencies),
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(float64(loops) / timeTaken.Seconds()),
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

	before := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < loops; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			beforeLocal := time.Now()
			_, errSel := status.SelectEntityoneOneByPK(dbConn, benchSQLLink, testEntityoneIDs[i%len(testEntityoneIDs)])
			if errSel != nil {
				errorC <- errSel
			} else {
				latenciesC <- time.Since(beforeLocal)
			}
		}(&wg)
	}

	var latencies []time.Duration
	var errCount int
	go receiveResults(&latenciesC, &errorC, &latencies, &errCount)

	wg.Wait()
	timeTaken := time.Since(before)

	return BenchResult{
			Action:     "selectEntityoneOneByPK",
			Loops:      loops,
			PauseTime:  0,
			Errors:     errCount,
			Min:        getMin(latencies),
			Max:        getMax(latencies),
			Median:     getMedian(latencies),
			StandDev:   getStandardDeviation(latencies),
			Throughput: int(float64(loops) / timeTaken.Seconds()),
		},
		nil
}

// receiveResults take the 2 channels used to receive results and gather the data into accessible variables
func receiveResults(
	latenciesC *chan time.Duration, errorC *chan error,
	lats *[]time.Duration, errCount *int,
) {
	for {
		select {
		case latency := <-*latenciesC:
			*lats = append(*lats, latency)
		case <-*errorC:
			*errCount++
		}
	}
}

// getMin retrieves the min latency
func getMin(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	return latencies[0]
}

// getMax retrieves the max latency
func getMax(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	return latencies[len(latencies)-1]
}

// getMedian returns the median duration of a list
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

// getStandardDeviation returns the standard deviation of the list
func getStandardDeviation(latencies []time.Duration) time.Duration {

	if len(latencies) == 0 {
		return 0
	}

	// Sum the square of the mean subtracted from each number
	mean := getMean(latencies)

	var variance float64

	for _, latency := range latencies {
		variance += math.Pow(float64(latency.Nanoseconds()-mean.Nanoseconds()), 2)
	}

	return time.Duration(math.Sqrt(variance / float64(len(latencies))))
}

// getMean returns the mean of the list
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
