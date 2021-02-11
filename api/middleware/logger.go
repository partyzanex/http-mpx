package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/partyzanex/http-mpx/api"
)

const logFormat = "method=%s uri=%s latency=%s"

// Logger middleware.
func Logger(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer elapsed(r)()
		return next(w, r)
	}
}

// elapsed calculated request execution time
// and write it in log.
func elapsed(r *http.Request) func() {
	start := time.Now()

	return func() {
		log.Printf(logFormat, r.Method, r.RequestURI, time.Since(start))
	}
}
