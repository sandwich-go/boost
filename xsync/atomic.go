package xsync

import (
	"sync/atomic"
	"time"
)

// AtomicBool is an atomic type-safe wrapper for bool values.
type AtomicBool struct{ v int32 }

// Set Store atomically stores the passed bool.
func (b *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreInt32(&b.v, 1)
	} else {
		atomic.StoreInt32(&b.v, 0)
	}
}

// Get Load atomically loads the wrapped bool.
func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32(&b.v) == 1
}

// AtomicInt32 is an atomic type-safe wrapper for int32 values.
type AtomicInt32 int32

// Add atomically adds n to *AtomicInt32 and returns the new value.
func (i *AtomicInt32) Add(n int32) int32 {
	return atomic.AddInt32((*int32)(i), n)
}

// Set Store atomically stores the passed int32.
func (i *AtomicInt32) Set(n int32) {
	atomic.StoreInt32((*int32)(i), n)
}

// Get Load atomically loads the wrapped int32.
func (i *AtomicInt32) Get() int32 {
	return atomic.LoadInt32((*int32)(i))
}

// CompareAndSwap executes the compare-and-swap operation for a int32 value.
func (i *AtomicInt32) CompareAndSwap(oldval, newval int32) (swapped bool) {
	return atomic.CompareAndSwapInt32((*int32)(i), oldval, newval)
}

// AtomicUint32 is an atomic type-safe wrapper for uint32 values.
type AtomicUint32 uint32

// Add atomically adds n to *AtomicUint32 and returns the new value.
func (i *AtomicUint32) Add(n uint32) uint32 {
	return atomic.AddUint32((*uint32)(i), n)
}

// Set Store atomically stores the passed uint32.
func (i *AtomicUint32) Set(n uint32) {
	atomic.StoreUint32((*uint32)(i), n)
}

// Get Load atomically loads the wrapped uint32.
func (i *AtomicUint32) Get() uint32 {
	return atomic.LoadUint32((*uint32)(i))
}

// CompareAndSwap executes the compare-and-swap operation for a uint32 value.
func (i *AtomicUint32) CompareAndSwap(oldval, newval uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32((*uint32)(i), oldval, newval)
}

// AtomicInt64 is an atomic type-safe wrapper for int64 values.
type AtomicInt64 int64

// Add atomically adds n to *AtomicInt64 and returns the new value.
func (i *AtomicInt64) Add(n int64) int64 {
	return atomic.AddInt64((*int64)(i), n)
}

// Set Store atomically stores the passed int64.
func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64((*int64)(i), n)
}

// Get Load atomically loads the wrapped int64.
func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

// CompareAndSwap executes the compare-and-swap operation for a int64 value.
func (i *AtomicInt64) CompareAndSwap(oldval, newval int64) (swapped bool) {
	return atomic.CompareAndSwapInt64((*int64)(i), oldval, newval)
}

// AtomicUint64 is an atomic type-safe wrapper for uint64 values.
type AtomicUint64 uint64

// Add atomically adds n to *AtomicUint64 and returns the new value.
func (i *AtomicUint64) Add(n uint64) uint64 {
	return atomic.AddUint64((*uint64)(i), n)
}

// Set Store atomically stores the passed uint64.
func (i *AtomicUint64) Set(n uint64) {
	atomic.StoreUint64((*uint64)(i), n)
}

// Get Load atomically loads the wrapped uint64.
func (i *AtomicUint64) Get() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}

// CompareAndSwap executes the compare-and-swap operation for a uint64 value.
func (i *AtomicUint64) CompareAndSwap(oldval, newval uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64((*uint64)(i), oldval, newval)
}

// AtomicDuration is an atomic type-safe wrapper for duration values.
type AtomicDuration int64

// Add atomically adds duration to *AtomicDuration and returns the new value.
func (d *AtomicDuration) Add(duration time.Duration) time.Duration {
	return time.Duration(atomic.AddInt64((*int64)(d), int64(duration)))
}

// Set Store atomically stores the passed time.Duration.
func (d *AtomicDuration) Set(duration time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(duration))
}

// Get Load atomically loads the wrapped time.Duration.
func (d *AtomicDuration) Get() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

// CompareAndSwap executes the compare-and-swap operation for an time.Duration value
func (d *AtomicDuration) CompareAndSwap(oldval, newval time.Duration) (swapped bool) {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(oldval), int64(newval))
}

// AtomicString is an atomic type-safe wrapper for string values.
type AtomicString struct {
	v atomic.Value
}

var _zeroString string

// NewAtomicString creates a new String.
func NewAtomicString(v string) *AtomicString {
	x := &AtomicString{}
	if v != _zeroString {
		x.Set(v)
	}
	return x
}

// Get Load atomically loads the wrapped string.
func (x *AtomicString) Get() string {
	if v := x.v.Load(); v != nil {
		return v.(string)
	}
	return _zeroString
}

// Set Store atomically stores the passed string.
func (x *AtomicString) Set(v string) {
	x.v.Store(v)
}

// AtomicTime is an atomic type-safe wrapper for time.Time values.
type AtomicTime struct {
	v int64
}

// Set Store atomically stores the passed time.Time.
func (t *AtomicTime) Set(v time.Time) {
	atomic.StoreInt64(&t.v, v.UnixNano())
}

// Get Load atomically loads the wrapped time.Time.
func (t *AtomicTime) Get() time.Time {
	return time.Unix(0, atomic.LoadInt64(&t.v))
}

// CompareAndSwap executes the compare-and-swap operation for a time.Time value.
func (t *AtomicTime) CompareAndSwap(oldval, newval time.Time) bool {
	return atomic.CompareAndSwapInt64(&t.v, oldval.UnixNano(), newval.UnixNano())
}
