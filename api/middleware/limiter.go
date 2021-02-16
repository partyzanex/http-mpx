package middleware

import (
	"net/http"

	"github.com/partyzanex/http-mpx/api"
	"github.com/partyzanex/http-mpx/pkg/limiter"
)

// ConcurrentLimiter creates api.Middleware function
// inside which the limiter is implemented.
func ConcurrentLimiter(limit int) api.Middleware {
	return func(next api.Handler) api.Handler {
		lm := limiter.Concurrent(limit)

		return func(w http.ResponseWriter, r *http.Request) (err error) {
			// apply limiter
			if !lm.Allow() {
				return api.NewError(http.StatusTooManyRequests, "")
			}

			// take a limit of calls
			release := lm.Take()
			defer release()

			err = next(w, r)

			return
		}
	}
}
