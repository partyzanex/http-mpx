package limiter_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/partyzanex/http-mpx/pkg/limiter"
)

func TestConcurrent_Take(t *testing.T) {
	limit := 1000
	lm := limiter.Concurrent(limit)
	counter := int32(0)

	fn := func() {
		release := lm.Take()
		defer release()

		atomic.AddInt32(&counter, 1)
		defer atomic.AddInt32(&counter, -1)

		if counter > int32(limit) {
			t.Error("limit exceeded")
		}
	}
	wg := sync.WaitGroup{}

	for i := 0; i < 100000; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			fn()
		}()
	}

	wg.Wait()
}
