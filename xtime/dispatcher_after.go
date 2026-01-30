//go:build !timewheel

package xtime

import "time"

var AfterFunc = time.AfterFunc
