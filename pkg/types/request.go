package types

import (
	"net/http"
)

type Requests []*Request

type Request struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}
