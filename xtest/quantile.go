package xtest

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

var durationsMap = make(map[string][]time.Duration)
var durationsMapMutex sync.RWMutex

// AddDuration 添加key所用的耗时
func AddDuration(key string, duration time.Duration) {
	durationsMapMutex.Lock()
	defer durationsMapMutex.Unlock()
	durations := durationsMap[key]
	durations = append(durations, duration)
	durationsMap[key] = durations
}

// PrintPercentiles 打印耗时数据
func PrintPercentiles(percentiles ...float32) {
	durationsMapMutex.Lock()
	defer durationsMapMutex.Unlock()
	results := CalculatePercentiles(durationsMap, percentiles)
	fmt.Printf("=== xtest package cost percentiles result(%d) === \n", len(durationsMap))
	for key := range durationsMap {
		fmt.Printf("Key: %s\n", key)
		for i, percentile := range percentiles {
			fmt.Printf("    Percentile %.2f: %.2f seconds\n", percentile, results[key][i].Seconds())
		}
		fmt.Println()
	}
}

// CalculatePercentiles 计算耗时数据
func CalculatePercentiles(durationsMap map[string][]time.Duration, percentiles []float32) map[string][]time.Duration {
	results := make(map[string][]time.Duration)
	for key, durations := range durationsMap {
		// Sort the durations in ascending order
		sortedDurations := make([]time.Duration, len(durations))
		copy(sortedDurations, durations)
		sort.Slice(sortedDurations, func(i, j int) bool {
			return sortedDurations[i] < sortedDurations[j]
		})
		// Calculate percentiles
		percentileDurations := make([]time.Duration, len(percentiles))
		for i, p := range percentiles {
			index := int(float32(len(sortedDurations)-1) * p)
			percentileDurations[i] = sortedDurations[index]
		}
		results[key] = percentileDurations
	}
	return results
}
