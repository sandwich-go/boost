package xsync

import (
	"context"
	"errors"
	"golang.org/x/sync/semaphore"
)

type RWTimeoutLock struct {
	sem *semaphore.Weighted
	r   int64 // max concurrent readers capacity
}

// RLock 在 ctx 取消或超时前获取读锁；失败返回 ctx.Err() 等。
func (l *RWTimeoutLock) RLock(ctx context.Context) error {
	return l.sem.Acquire(ctx, 1)
}

func (l *RWTimeoutLock) RUnlock() {
	l.sem.Release(1)
}

func (l *RWTimeoutLock) Lock(ctx context.Context) error {
	return l.sem.Acquire(ctx, l.r)
}

func (l *RWTimeoutLock) Unlock() {
	l.sem.Release(l.r)
}

func NewRWTimeoutLock(maxReaders int64) *RWTimeoutLock {
	if maxReaders <= 0 {
		panic(errors.New("maxReaders must be greater than zero"))
	}
	return &RWTimeoutLock{
		sem: semaphore.NewWeighted(maxReaders),
		r:   maxReaders,
	}
}
