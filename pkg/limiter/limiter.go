package limiter

import (
	"sync/atomic"
	"time"
)

type rateLimiter struct {
	d       time.Duration
	limit   int64
	counter *int64
}

func New(d time.Duration, limit int64) Limiter {
	counter := int64(0)
	lim := rateLimiter{
		d:       d,
		limit:   limit,
		counter: &counter,
	}

	go lim.start()

	return &lim
}

func (lim *rateLimiter) start() {
	for {
		select {
		case <-time.After(lim.d):
			atomic.AddInt64(lim.counter, -*lim.counter)
		}
	}
}

func (lim *rateLimiter) Allow() bool {
	if *lim.counter < lim.limit {
		atomic.AddInt64(lim.counter, 1)
		return true
	}

	return false
}
