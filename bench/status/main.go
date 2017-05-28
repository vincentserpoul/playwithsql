package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"

	"time"

	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/bench"
	"github.com/vincentserpoul/playwithsql/dbhandler"
	"github.com/vincentserpoul/playwithsql/status"
)

// Number of retries after query returns an error
const maxRetryCount = 3

// Results to be returned
type Results struct {
	DBType       string
	MaxConns     int
	Date         time.Time
	BenchResults []BenchResult
}

// BenchResult data
type BenchResult struct {
	Action      string
	Loops       int
	Concurrency int
	PauseTime   time.Duration
	Errors      int
	Min         time.Duration
	Max         time.Duration
	Median      time.Duration
	StandDev    time.Duration
	Throughput  int
}

func main() {

	// Flags
	dbName := "playwithsql"
	dbType := flag.String("db", "mysql", "type of db to bench: mysql, cockroachdb, postgres...")
	schemaType := flag.String("sch", "islatest", "type of schema to use, is latest, latest status...")
	dbHost := flag.String("host", "127.0.0.1", "host IP")
	loops := flag.Int("loops", 100, "number of loops")
	maxConns := flag.Int("maxconns", 10, "number of max connections")
	concurrency := flag.Int("concurrency", 100, "number of go routines created")

	flag.Parse()

	if *concurrency > *loops {
		*concurrency = *loops
	}

	db, err := dbhandler.Get(*dbType, *dbHost, dbName)
	if err != nil {
		log.Fatalf("%s - %s - %s, \n%v", *dbType, *dbHost, dbName, err)
	}

	// Connection
	islatestSQLLink := status.GetSQLIntImpl(*dbType, *schemaType)

	ctx := context.Background()
	err = islatestSQLLink.MigrateDown(ctx, db)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = islatestSQLLink.MigrateUp(ctx, db)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Number of max connections
	// TODO set the param in the db config
	db.SetMaxOpenConns(*maxConns)
	db.SetMaxIdleConns(*maxConns)

	var results = Results{
		DBType:   *dbType,
		MaxConns: *maxConns,
		Date:     time.Now(),
	}

	// Create
	createResults, testEntityoneIDs, err := BenchmarkCreate(ctx, *loops, *concurrency, db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, createResults)

	// Update
	updateResults, err := BenchmarkUpdateStatus(ctx, *loops, *concurrency, db, islatestSQLLink, testEntityoneIDs)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, updateResults)

	// Select by status
	selectByStatusResults, err := BenchmarkSelectEntityoneByStatus(ctx, *loops, *concurrency, db, islatestSQLLink)
	if err != nil {
		log.Fatalf("%v", err)
	}
	results.BenchResults = append(results.BenchResults, selectByStatusResults)

	// Select by PK
	selectByPKResults, err := BenchmarkSelectEntityoneOneByPK(ctx, *loops, *concurrency, db, islatestSQLLink, testEntityoneIDs)
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
	ctx context.Context,
	loops int,
	concurrency int,
	dbConn *sqlx.DB,
	benchSQLLink *status.SQLIntImpl,
) (
	results BenchResult,
	testEntityoneIDs []int64,
	err error,
) {
	entityIDsC := make(chan int64)
	latenciesC, errorC, doneC, resultsC := handleResults()

	before := time.Now()
	var wg sync.WaitGroup

	// Pause time
	dynPauseTime := 1 * time.Millisecond
	dynPauseTimeC := dynPauseTimeInit(&dynPauseTime)
	defer close(dynPauseTimeC)

	// Launch as many concurrent connections as asked
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(routineNum int, ctx context.Context, wg *sync.WaitGroup) {
			for j := 0; j < loops/concurrency; j++ {
				time.Sleep(dynPauseTime)
				var e status.Entityone
				beforeLocal := time.Now()
				ok := false
				var errCr error
				retryCount := 0
				for retryCount < maxRetryCount && !ok {
					// Timeout
					sqlCtx, sqlCncl := context.WithTimeout(ctx, 250*time.Millisecond)
					defer sqlCncl()

					// For each error, we add some pause time
					errCr = e.Create(sqlCtx, dbConn, benchSQLLink)
					if errCr != nil {
						retryCount++
						time.Sleep(dynPauseTime)
						dynPauseTimeC <- 1 * time.Millisecond
					} else {
						ok = true
					}
				}
				if errCr != nil {
					errorC <- errCr
				} else {
					latenciesC <- time.Since(beforeLocal)
					entityIDsC <- e.ID
					// If no error, we increment down a little bit
					dynPauseTimeC <- -1 * time.Millisecond
				}
			}
			wg.Done()
		}(i, ctx, &wg)
	}

	// Receive the entityIDs
	go func() {
		for entityID := range entityIDsC {
			testEntityoneIDs = append(testEntityoneIDs, entityID)
		}
	}()

	wg.Wait()
	doneC <- true
	rawRes := <-resultsC
	timeTaken := time.Since(before)
	return BenchResult{
			Action:      "create",
			Loops:       loops,
			Concurrency: concurrency,
			PauseTime:   dynPauseTime,
			Errors:      rawRes.errCount,
			Min:         bench.GetMin(rawRes.latencies),
			Max:         bench.GetMax(rawRes.latencies),
			Median:      bench.GetMedian(rawRes.latencies),
			StandDev:    bench.GetStandardDeviation(rawRes.latencies),
			Throughput:  int(float64(loops) / timeTaken.Seconds()),
		},
		testEntityoneIDs,
		nil
}

// BenchmarkUpdateStatus benchmark for status updates (include deletes)
func BenchmarkUpdateStatus(
	ctx context.Context,
	loops int,
	concurrency int,
	dbConn *sqlx.DB,
	benchSQLLink *status.SQLIntImpl,
	testEntityoneIDs []int64,
) (
	results BenchResult,
	err error,
) {
	if len(testEntityoneIDs) == 0 {
		return results, fmt.Errorf("BenchmarkUpdateStatus: no entity created, nothing to update")
	}

	latenciesC, errorC, doneC, resultsC := handleResults()

	before := time.Now()
	var wg sync.WaitGroup

	// Pause time
	dynPauseTime := 1 * time.Millisecond
	dynPauseTimeC := dynPauseTimeInit(&dynPauseTime)
	defer close(dynPauseTimeC)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func(ctx context.Context, wg *sync.WaitGroup) {
			for j := 0; j < loops/concurrency; j++ {
				time.Sleep(dynPauseTime)
				var e status.Entityone
				e.ID = testEntityoneIDs[i%len(testEntityoneIDs)]
				beforeLocal := time.Now()
				ok := false
				var errU error
				retryCount := 0
				for retryCount < maxRetryCount && !ok {
					// Timeout
					sqlCtx, sqlCncl := context.WithTimeout(ctx, 250*time.Millisecond)
					defer sqlCncl()
					errU = e.UpdateStatus(sqlCtx, dbConn, benchSQLLink, status.ActionCancel, status.StatusCancelled)
					if errU != nil {
						retryCount++
						time.Sleep(dynPauseTime)
						dynPauseTimeC <- 1 * time.Millisecond
					} else {
						ok = true
					}
				}
				if errU != nil {
					errorC <- errU
				} else {
					latenciesC <- time.Since(beforeLocal)
					// If no error, we increment down a little bit
					dynPauseTimeC <- -1 * time.Millisecond
				}
			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
	doneC <- true
	rawRes := <-resultsC
	timeTaken := time.Since(before)

	return BenchResult{
			Action:      "updateStatus",
			Loops:       loops,
			Concurrency: concurrency,
			PauseTime:   dynPauseTime,
			Errors:      rawRes.errCount,
			Min:         bench.GetMin(rawRes.latencies),
			Max:         bench.GetMax(rawRes.latencies),
			Median:      bench.GetMedian(rawRes.latencies),
			StandDev:    bench.GetStandardDeviation(rawRes.latencies),
			Throughput:  int(float64(loops) / timeTaken.Seconds()),
		},
		nil

}

// BenchmarkSelectEntityoneByStatus benchmark with select by status
func BenchmarkSelectEntityoneByStatus(
	ctx context.Context,
	loops int,
	concurrency int,
	dbConn *sqlx.DB,
	benchSQLLink *status.SQLIntImpl,
) (
	results BenchResult,
	err error,
) {
	latenciesC, errorC, doneC, resultsC := handleResults()

	var wg sync.WaitGroup
	before := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup) {
			for j := 0; j < loops/concurrency; j++ {
				beforeLocal := time.Now()
				sqlCtx, sqlCncl := context.WithTimeout(ctx, 100*time.Millisecond)
				defer sqlCncl()
				_, errSel := status.SelectEntityoneByStatus(sqlCtx, dbConn, benchSQLLink, status.StatusCancelled)
				if errSel != nil {
					errorC <- errSel
				} else {
					latenciesC <- time.Since(beforeLocal)
				}
			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
	doneC <- true
	rawRes := <-resultsC
	timeTaken := time.Since(before)

	return BenchResult{
			Action:      "selectEntityoneByStatus",
			Loops:       loops,
			Concurrency: concurrency,
			PauseTime:   0,
			Errors:      rawRes.errCount,
			Min:         bench.GetMin(rawRes.latencies),
			Max:         bench.GetMax(rawRes.latencies),
			Median:      bench.GetMedian(rawRes.latencies),
			StandDev:    bench.GetStandardDeviation(rawRes.latencies),
			Throughput:  int(float64(loops) / timeTaken.Seconds()),
		},
		nil
}

// BenchmarkSelectEntityoneOneByPK benchmark with select by primary key
func BenchmarkSelectEntityoneOneByPK(
	ctx context.Context,
	loops int,
	concurrency int,
	dbConn *sqlx.DB,
	benchSQLLink *status.SQLIntImpl,
	testEntityoneIDs []int64,
) (
	results BenchResult,
	err error,
) {
	latenciesC, errorC, doneC, resultsC := handleResults()

	before := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < loops; i++ {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup) {
			for j := 0; j < loops/concurrency; j++ {
				beforeLocal := time.Now()
				sqlCtx, sqlCncl := context.WithTimeout(ctx, 100*time.Millisecond)
				defer sqlCncl()
				_, errSel := status.SelectEntityoneOneByPK(sqlCtx, dbConn, benchSQLLink, testEntityoneIDs[i%len(testEntityoneIDs)])
				if errSel != nil {
					errorC <- errSel
				} else {
					latenciesC <- time.Since(beforeLocal)
				}
			}
			wg.Done()
		}(ctx, &wg)
	}

	wg.Wait()
	doneC <- true
	rawRes := <-resultsC
	timeTaken := time.Since(before)

	return BenchResult{
			Action:      "selectEntityoneOneByPK",
			Loops:       loops,
			Concurrency: concurrency,
			PauseTime:   0,
			Errors:      rawRes.errCount,
			Min:         bench.GetMin(rawRes.latencies),
			Max:         bench.GetMax(rawRes.latencies),
			Median:      bench.GetMedian(rawRes.latencies),
			StandDev:    bench.GetStandardDeviation(rawRes.latencies),
			Throughput:  int(float64(loops) / timeTaken.Seconds()),
		},
		nil
}

type rawResults struct {
	latencies []time.Duration
	errCount  int
}

// handleResults will generate two channels that will receive latencies and errors
func handleResults() (chan time.Duration, chan error, chan bool, chan rawResults) {
	latenciesC := make(chan time.Duration)
	errorC := make(chan error)
	resultsC := make(chan rawResults)
	doneC := make(chan bool)

	var latencies []time.Duration
	var errCount int
	var mux sync.Mutex

	go func() {
		for {
			select {
			case latency := <-latenciesC:
				latencies = append(latencies, latency)
			case erRrrR := <-errorC:
				fmt.Println(erRrrR)
				mux.Lock()
				errCount++
				mux.Unlock()
			case <-doneC:
				resultsC <- rawResults{latencies: latencies, errCount: errCount}
			}
		}
	}()

	return latenciesC, errorC, doneC, resultsC
}

const (
	maxPauseTime = 200 * time.Millisecond
	minPauseTime = 1 * time.Millisecond
)

// dynPauseTimeInit generates a channel that will be used to dynamically update the pause time between transactions
func dynPauseTimeInit(dynPauseTime *time.Duration) chan time.Duration {
	dynPauseTimeC := make(chan time.Duration)
	go func() {
		for additionalPauseTime := range dynPauseTimeC {
			if (*dynPauseTime+additionalPauseTime) > minPauseTime && (*dynPauseTime+additionalPauseTime) < maxPauseTime {
				*dynPauseTime += additionalPauseTime
			}
		}
	}()
	return dynPauseTimeC
}
