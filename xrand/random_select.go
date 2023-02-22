package xrand

import (
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type int32Slice []int32

func (p int32Slice) Len() int           { return len(p) }
func (p int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// RandomInt 随机某个值，该值[min, max]
func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomSelectOneFromMap map[int32]int32{1:10,2:24,3:53}
func RandomSelectOneFromMap(numbers map[int32]int32) (int32, bool) {
	var total, tmpTotal int32
	keyArr := make([]int32, 0, len(numbers))
	for key, rate := range numbers {
		total += rate
		keyArr = append(keyArr, key)
	}
	if total <= 0 {
		return 0, false
	}
	sort.Sort(int32Slice(keyArr))
	r := rand.Int31n(total) + 1
	for _, key := range keyArr {
		rate := numbers[key]
		tmpTotal += rate
		if r <= tmpTotal {
			return key, true
		}
	}
	return 0, false
}

// RandomSelectOneFromArray array(20,30,50) 返回0的概率0.2 返回1的概率0.3 返回2的概率0.5
func RandomSelectOneFromArray(numbers []int32) (int, bool) {
	var total, tmpTotal int32
	for _, rate := range numbers {
		total += rate
	}
	if total <= 0 {
		return 0, false
	}
	r := rand.Int31n(total) + 1
	for index, rate := range numbers {
		tmpTotal += rate
		if r <= tmpTotal {
			return index, true
		}
	}
	return 0, false
}

// IsSelected100n  number / 100 的概率返回true
func IsSelected100n(number int32) bool {
	return rand.Int31n(100)+1 <= number
}
