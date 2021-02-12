package pool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/partyzanex/http-mpx/internal/assert"
	"github.com/partyzanex/http-mpx/pkg/pool"
)

func TestNewWorkersPool(t *testing.T) {
	res := int64(0)
	wp := pool.New(10, pool.WithBuffer(1000))

	for i := 0; i < 1000; i++ {
		wp.Add(func() {
			atomic.AddInt64(&res, 1)
		})
	}

	wp.Wait()

	assert.Equal(t, int64(1000), res)
}

func TestWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := pool.New(10, pool.WithContext(ctx))
	res := int64(0)

	for i := 0; i < 100; i++ {
		wp.Add(func() {
			if res >= 50 {
				cancel()
			}

			atomic.AddInt64(&res, 1)
		})
	}

	wp.Wait()

	assert.Equal(t, true, res < 100)
}

func BenchmarkPool_Add1(b *testing.B) {
	fn := func() {}

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		wp := pool.New(10)

		b.StartTimer()

		wp.Add(fn)

		b.StopTimer()

		wp.Wait()
	}
}

func BenchmarkPool_Add10(b *testing.B) {
	fn := func() {}

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		wp := pool.New(10, pool.WithBuffer(10))

		b.StartTimer()

		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)
		wp.Add(fn)

		b.StopTimer()

		wp.Wait()
	}
}

func BenchmarkNew(b *testing.B) {
	buf := pool.WithBuffer(10)

	for i := 0; i < b.N; i++ {
		p := pool.New(100, buf)
		p.Wait()
	}
}

func BenchmarkWaitGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		wg.Add(10)

		fn := func() {
			wg.Done()
		}
		fn()
		fn()
		fn()
		fn()
		fn()
		fn()
		fn()
		fn()
		fn()
		fn()
		wg.Wait()
	}
}
