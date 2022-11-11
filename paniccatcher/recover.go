package paniccatcher

import (
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
	"time"
)

func AutoRecover(tag string, f func(), opts ...AutoRecoverOption) {
	defer func() {
		if reason := recover(); reason != nil {
			cc := NewAutoRecoverOptions(opts...)
			if cc.OnRecover != nil {
				cc.OnRecover(tag, reason)
			} else {
				log.Error(fmt.Sprintf("%s panic with err, reason: %v", tag, reason))
			}
			if cc.DelayTime > 0 {
				time.Sleep(cc.DelayTime)
			}
			go AutoRecover(tag, f, opts...)
		}
	}()
	f()
}
