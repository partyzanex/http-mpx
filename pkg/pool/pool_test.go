package pool_test

import (
	"sync"
	"testing"

	"github.com/partyzanex/http-mpx/pkg/pool"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkersPool(t *testing.T) {
	wp := pool.New(10)
	res := make([]int, 0)
	mu := sync.Mutex{}

	for i := 0; i < 100; i++ {
		wp.Add(func() {
			mu.Lock()
			res = append(res, 1)
			mu.Unlock()
		})
	}

	wp.Wait()

	assert.Equal(t, 100, len(res))
}
