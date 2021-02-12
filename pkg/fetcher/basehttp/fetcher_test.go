package basehttp_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/partyzanex/http-mpx/internal/assert"
	"github.com/partyzanex/http-mpx/pkg/fetcher/basehttp"
	"github.com/partyzanex/http-mpx/pkg/types"
)

func TestHttpFetcher_Fetch(t *testing.T) {
	req := types.Request{
		Method: http.MethodPatch,
		Headers: http.Header{
			"Test-Header":     []string{"test header value"},
			"Content-Type":    []string{"text/plain; charset=utf-8"},
			"Accept-Encoding": []string{"gzip"},
			"Content-Length":  []string{"9"},
			"User-Agent":      []string{"Go-http-client/1.1"},
		},
		Body: []byte("test body"),
	}
	f := basehttp.New(nil)
	expBody := []byte("test response")
	expStatus := http.StatusGone

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, req.Method, r.Method)
		assert.Equal(t, req.Headers, r.Header)

		body, err := ioutil.ReadAll(r.Body)
		assert.Equal(t, nil, err)
		assert.Equal(t, req.Body, body)

		w.WriteHeader(http.StatusGone)

		_, err = w.Write([]byte("test response"))
		assert.Equal(t, nil, err)
	}))
	defer server.Close()

	req.URL = server.URL

	result, err := f.Fetch(context.Background(), req)
	assert.Equal(t, nil, err)

	assert.Equal(t, expBody, result.Body)
	assert.Equal(t, expStatus, result.StatusCode)
	assert.Equal(t, server.URL, result.URL)

	req.URL = ""
	req.Method = ""

	var ctx context.Context

	_, err = f.Fetch(ctx, req)
	assert.Equal(t, false, err == nil)

	_, err = f.Fetch(context.Background(), req)
	assert.Equal(t, false, err == nil)
}
