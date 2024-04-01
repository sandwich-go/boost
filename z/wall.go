package z

import "time"

type wallClock struct {
	t      time.Time
	offset time.Duration
}

var wc wallClock

func init() {
	wc = wallClock{t: time.Now(), offset: offset()}
}

// Now returns the current wall clock time.
func Now() time.Time {
	return wc.t.Add(offset() - wc.offset)
}

// NowWithOffset returns the wall clock time with given offset.
func NowWithOffset(monoOffset time.Duration) time.Time {
	return wc.t.Add(monoOffset - wc.offset)
}

func offset() time.Duration {
	return time.Duration(runtimeNanotime())
}
