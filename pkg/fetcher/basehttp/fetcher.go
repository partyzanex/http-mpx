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

// httpFetcher implements of fetcher.Interface.
type httpFetcher struct {
	*http.Client
}

// New creates a fetcher.Interface instance.
func New(client *http.Client) fetcher.Interface {
	if client == nil {
		client = http.DefaultClient
	}

	return &httpFetcher{client}
}

// Fetch takes a types.Request, execute HTTP-request and returns *types.Result or error.
func (c *httpFetcher) Fetch(ctx context.Context, req types.Request) (*types.Result, error) {
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

	defer func() {
		erc := resp.Body.Close()

		if err == nil {
			err = erc
		}
	}()

	return c.handleResponse(resp)
}

// handleResponse parses *http.Response body and returns *types.Result or error.
func (*httpFetcher) handleResponse(resp *http.Response) (*types.Result, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %s", err)
	}

	return &types.Result{
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}
