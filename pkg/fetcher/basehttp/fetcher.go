package basehttp

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/partyzanex/http-mpx/pkg/fetcher"
	"github.com/partyzanex/http-mpx/pkg/types"
)

type httpFetcher struct {
	*http.Client
}

func New(client *http.Client) fetcher.Interface {
	if client == nil {
		client = http.DefaultClient
	}

	return &httpFetcher{
		Client: client,
	}
}

func (c *httpFetcher) Fetch(ctx context.Context, req types.Request) (result *types.Result, err error) {
	method := req.Method

	if method == "" {
		method = http.MethodGet
	}

	r, err := http.NewRequestWithContext(ctx, method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %s", err)
	}

	r.Header = req.Headers

	resp, err := c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("cannot load response: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %s", err)
	}

	result = &types.Result{
		URL:        req.URL,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}

	return
}
