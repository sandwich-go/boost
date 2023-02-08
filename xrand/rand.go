package xrand

import "github.com/sandwich-go/boost/z"

var (
	// FastRand is a fast thread local random function
	FastRand = z.FastRand
	// FastRandUint32n returns pseudorandom uint32 in the range [0..maxN)
	FastRandUint32n = z.FastRandUint32n
)
