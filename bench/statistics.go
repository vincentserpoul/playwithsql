package bench

import (
	"math"
	"sort"
	"time"
)

// GetMin retrieves the min latency
func GetMin(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	return latencies[0]
}

// GetMax retrieves the max latency
func GetMax(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	return latencies[len(latencies)-1]
}

// GetMedian returns the median duration of a list
func GetMedian(latencies []time.Duration) time.Duration {
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

// GetStandardDeviation returns the standard deviation of the list
func GetStandardDeviation(latencies []time.Duration) time.Duration {

	if len(latencies) == 0 {
		return 0
	}

	// Sum the square of the mean subtracted from each number
	mean := GetMean(latencies)

	var variance float64

	for _, latency := range latencies {
		variance += math.Pow(float64(latency.Nanoseconds()-mean.Nanoseconds()), 2)
	}

	return time.Duration(math.Sqrt(variance / float64(len(latencies))))
}

// GetMean returns the mean of the list
func GetMean(latencies []time.Duration) time.Duration {
	if len(latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, latency := range latencies {
		total += latency
	}

	return time.Duration(total.Nanoseconds() / int64(len(latencies)))
}
