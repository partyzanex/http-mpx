package types

import "net/http"

// Result represents a result item that is to be returned as a response.
type Result struct {
	// requested URL
	URL string `json:"url"`
	// HTTP Status Code number
	StatusCode int `json:"status_code"`
	// Returned HTTP headers
	Headers http.Header `json:"headers"`
	// Returned content
	Body []byte `json:"body"`
}
