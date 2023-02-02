package paniccatcher

import (
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
	"time"
)

// AutoRecover 执行 f ，当 f 发生panic，启动新的协程重复 AutoRecover.
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
