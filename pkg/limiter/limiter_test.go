package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/partyzanex/http-mpx/pkg/limiter"
)

func TestRateLimiter_Allow(t *testing.T) {
	rate := limiter.New(100*time.Millisecond, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	i := 0

	for {
		select {
		case <-ctx.Done():
			assert.Equal(t, 100, i)
			return
		default:
			allow := rate.Allow()
			if allow {
				i++
			}

			if !allow {
				cancel()
			}
		}
	}
}
