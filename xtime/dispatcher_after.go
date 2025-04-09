//go:build !timingwheel

package xtime

import "time"

var timeAfterFunc = time.AfterFunc

type timer *time.Timer
