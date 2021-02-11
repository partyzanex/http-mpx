package basehttp_test

import (
	"context"
	"net/http"
	"net/url"
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

	uri, err := url.Parse("https://ya.ru/")
	assert.Equal(t, nil, err)

	req := types.Request{
		URL:     uri,
		Method:  http.MethodGet,
		Headers: nil,
		Body:    nil,
	}

	result, err := f.Fetch(ctx, req)
	assert.Equal(t, nil, err)

	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.Equal(t, req.URL.String(), result.URL)
	assert.Equal(t, true, len(result.Headers) > 0)
	assert.Equal(t, true, len(result.Body) > 0)
}
