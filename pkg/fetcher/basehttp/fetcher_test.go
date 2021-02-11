package basehttp_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/partyzanex/http-mpx/pkg/fetcher/basehttp"
	"github.com/partyzanex/http-mpx/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestHttpFetcher_Fetch(t *testing.T) {
	f := basehttp.New(nil)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)

	defer cancel()

	req := types.Request{
		URL:     "https://ya.ru/",
		Method:  http.MethodGet,
		Headers: nil,
		Body:    nil,
	}

	result, err := f.Fetch(ctx, req)
	assert.Equal(t, nil, err)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, req.URL, result.URL)
	assert.Equal(t, true, len(result.Headers) > 0)
	assert.Equal(t, true, len(result.Body) > 0)
}
