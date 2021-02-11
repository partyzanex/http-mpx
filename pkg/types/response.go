package types

import "net/http"

type Response map[string]*Result

type Result struct {
	URL        string      `json:"url"`
	StatusCode int         `json:"status_code"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
}
