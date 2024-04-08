package ratelimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	rl := New(100) // per second, some slack.

	rl.Take()                         // Initialize.
	time.Sleep(time.Millisecond * 45) // Let some time pass.

	prev := time.Now()
	for i := 0; i < 10; i++ {
		now := rl.Take()
		if i > 0 {
			fmt.Println(i, now.Sub(prev).Round(time.Millisecond*2))
		}
		prev = now
	}
}
