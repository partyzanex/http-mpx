package fetch_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/partyzanex/http-mpx/api"
	"github.com/partyzanex/http-mpx/api/fetch"
	"github.com/partyzanex/http-mpx/api/middleware"
	"github.com/partyzanex/http-mpx/internal/assert"
	"github.com/partyzanex/http-mpx/pkg/fetcher"
	"github.com/partyzanex/http-mpx/pkg/types"
)

func TestHandler(t *testing.T) {
	request := makeRequest()
	input := makeInput(request, 20)
	f := makeFetcher(t, request)
	handler := fetch.GetHandler(makeConfig(), f)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(toJSON(t, input)))
	rec := httptest.NewRecorder()

	err := handler(rec, req)
	assert.Equal(t, nil, err)

	if assert.Equal(t, http.StatusOK, rec.Code) {
		var results []*types.Result

		fromJSON(t, rec.Body.Bytes(), &results)

		for _, result := range results {
			assert.Equal(t, request.URL, result.URL)
			assert.Equal(t, request.Body, result.Body)
			assert.Equal(t, request.Headers, result.Headers)
			assert.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
		}
	}
}

func TestHandlerErrors(t *testing.T) {
	request := makeRequest()
	f := makeFetcher(t, request)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{a:b}")))
	rec := httptest.NewRecorder()
	handler := fetch.GetHandler(makeConfig(), f)
	handler = middleware.ErrorHandler(handler)

	err := handler(rec, req)
	assert.Equal(t, nil, err)

	if assert.Equal(t, http.StatusBadRequest, rec.Code) {
		apiErr := api.Error{}
		fromJSON(t, rec.Body.Bytes(), &apiErr)
		assert.Equal(t, fetch.MsgDecodeRequest, apiErr.Message)
	}

	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(toJSON(t, makeInput(request, 25))))
	rec = httptest.NewRecorder()

	err = handler(rec, req)
	assert.Equal(t, nil, err)

	if assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code) {
		apiErr := api.Error{}
		fromJSON(t, rec.Body.Bytes(), &apiErr)
		assert.Equal(t, fetch.MsgLargeRequest, apiErr.Message)
	}

	expErrMsg := "test error message"

	f.FetchMethod = func(ctx context.Context, got types.Request) (*types.Result, error) {
		return nil, errors.New(expErrMsg)
	}

	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(toJSON(t, makeInput(request, 20))))
	rec = httptest.NewRecorder()

	err = handler(rec, req)
	assert.Equal(t, nil, err)

	if assert.Equal(t, http.StatusExpectationFailed, rec.Code) {
		apiErr := api.Error{}
		fromJSON(t, rec.Body.Bytes(), &apiErr)
		assert.Equal(t, true, strings.Count(apiErr.Message, expErrMsg) > 0)
	}
}

func TestHandlerTimeout(t *testing.T) {
	request := makeRequest()
	f := makeFetcher(t, request)
	config := makeConfig()
	handler := fetch.GetHandler(config, f)
	input := makeInput(request, 19)
	expErrMsg := "timeout failed"
	errSent := false

	f.FetchMethod = func(ctx context.Context, got types.Request) (*types.Result, error) {
		if !errSent {
			select {
			case <-ctx.Done():
				break
			case <-time.After(time.Millisecond):
				errSent = true
				return nil, errors.New(expErrMsg)
			}
		}

		return nil, errors.New("fetch failed")
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(toJSON(t, input)))
	rec := httptest.NewRecorder()
	handler = middleware.ErrorHandler(handler)

	err := handler(rec, req)
	assert.Equal(t, nil, err)

	if assert.Equal(t, http.StatusExpectationFailed, rec.Code) {
		apiErr := api.Error{}
		fromJSON(t, rec.Body.Bytes(), &apiErr)
		assert.Equal(t, true, strings.Count(apiErr.Message, expErrMsg) > 0)
	}
}

func TestHandlerOutgoing(t *testing.T) {
	request := makeRequest()
	f := makeFetcher(t, request)
	config := makeConfig()
	handler := fetch.GetHandler(config, f)
	input := makeInput(request, 20)
	counter := int32(0)

	f.FetchMethod = func(ctx context.Context, got types.Request) (*types.Result, error) {
		atomic.AddInt32(&counter, 1)
		defer atomic.AddInt32(&counter, -1)

		if counter > int32(config.Outgoing) {
			return nil, errors.New("outgoing failed")
		}

		return nil, nil
	}

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(toJSON(t, input)))
	rec := httptest.NewRecorder()
	handler = middleware.ErrorHandler(handler)

	err := handler(rec, req)
	assert.Equal(t, nil, err)

	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Log(rec.Body.String())
	}
}

func makeRequest() types.Request {
	return types.Request{
		URL:    "test-url",
		Method: http.MethodDelete,
		Headers: http.Header{
			"Test-Header": []string{"test header value"},
		},
		Body: []byte("test body"),
	}
}

func makeInput(request types.Request, count int) []*types.Request {
	input := make([]*types.Request, count)

	for i := 0; i < count; i++ {
		input[i] = &request
	}

	return input
}

func makeFetcher(t testing.TB, request types.Request) *fetcher.Mockup {
	return &fetcher.Mockup{
		FetchMethod: func(ctx context.Context, got types.Request) (*types.Result, error) {
			assert.Equal(t, request, got)

			return &types.Result{
				URL:        got.URL,
				StatusCode: http.StatusMethodNotAllowed,
				Headers:    got.Headers,
				Body:       got.Body,
			}, nil
		},
	}
}

func makeConfig() fetch.Config {
	return fetch.Config{
		Outgoing: 4,
		MaxURls:  20,
		Timeout:  10 * time.Millisecond,
	}
}

func toJSON(t testing.TB, v interface{}) []byte {
	b, err := json.Marshal(v)
	assert.Equal(t, nil, err)

	return b
}

func fromJSON(t testing.TB, in []byte, v interface{}) {
	err := json.Unmarshal(in, v)
	assert.Equal(t, nil, err)
}

func BenchmarkHandler(b *testing.B) {
	request := makeRequest()
	f := makeFetcher(b, request)
	f.FetchMethod = func(_ context.Context, _ types.Request) (*types.Result, error) {
		return nil, nil
	}
	handler := fetch.GetHandler(makeConfig(), f)
	body := bytes.NewReader(toJSON(b, makeInput(request, 20)))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		req := httptest.NewRequest(http.MethodPost, "/", body)
		rec := httptest.NewRecorder()

		b.StartTimer()

		_ = handler(rec, req)
	}
}

func BenchmarkHandlerParallel(b *testing.B) {
	request := makeRequest()
	f := makeFetcher(b, request)
	f.FetchMethod = func(_ context.Context, _ types.Request) (*types.Result, error) {
		return nil, nil
	}
	handler := fetch.GetHandler(makeConfig(), f)
	body := bytes.NewReader(toJSON(b, makeInput(request, 20)))

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()

			req := httptest.NewRequest(http.MethodPost, "/", body)
			rec := httptest.NewRecorder()

			b.StartTimer()

			_ = handler(rec, req)
		}
	})
}
