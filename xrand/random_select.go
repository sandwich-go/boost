package xrand

import (
	"cmp"
	"github.com/sandwich-go/boost/misc"
	"golang.org/x/exp/constraints"
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// RandomInt 随机某个值，该值[min, max]
func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomSelectOneFromMap map[int32]int32{1:10,2:24,3:53}
func RandomSelectOneFromMap(numbers map[int32]int32) (int32, bool) {
	return RandomSelectOneKeyFromMap[int32, int32](numbers)
}

// RandomSelectOneKeyFromMap map[int32]int32{1:10,2:24,3:53}
func RandomSelectOneKeyFromMap[K cmp.Ordered, V constraints.Integer](numbers map[K]V) (K, bool) {
	var total, tmpTotal V
	keyArr := make([]K, 0, len(numbers))
	for key, rate := range numbers {
		total += rate
		keyArr = append(keyArr, key)
	}
	if total <= 0 {
		return misc.Zero[K](), false
	}
	sort.Slice(keyArr, func(i, j int) bool {
		return cmp.Compare(keyArr[i], keyArr[j]) < 0
	})
	r := rand.Int63n(int64(total)) + 1
	for _, key := range keyArr {
		rate := numbers[key]
		tmpTotal += rate
		if r <= int64(tmpTotal) {
			return key, true
		}
	}
	return misc.Zero[K](), false
}

// RandomSelectOneFromArray array(20,30,50) 返回0的概率0.2 返回1的概率0.3 返回2的概率0.5
func RandomSelectOneFromArray(numbers []int32) (int, bool) {
	return RandomSelectIndexFromArray[int32](numbers)
}

// RandomSelectIndexFromArray array(20,30,50) 返回0的概率0.2 返回1的概率0.3 返回2的概率0.5
func RandomSelectIndexFromArray[T constraints.Integer](numbers []T) (int, bool) {
	var total, tmpTotal T
	for _, rate := range numbers {
		total += rate
	}
	if total <= 0 {
		return 0, false
	}
	r := rand.Int63n(int64(total)) + 1
	for index, rate := range numbers {
		tmpTotal += rate
		if r <= int64(tmpTotal) {
			return index, true
		}
	}
	return 0, false
}

// IsSelected100n  number / 100 的概率返回true
func IsSelected100n(number int32) bool {
	return rand.Int31n(100)+1 <= number
}
