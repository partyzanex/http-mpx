package types

import (
	"net/http"
)

// Request represents a minimal request item for each URL
// that can be requested.
type Request struct {
	// requested URL
	URL string `json:"url"`
	// HTTP-method
	Method string `json:"method"`
	// to be sent HTTP-headers
	Headers http.Header `json:"headers"`
	// to be sent HTTP-body
	Body []byte `json:"body"`
}
