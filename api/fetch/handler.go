package fetch

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/partyzanex/http-mpx/api"
	"github.com/partyzanex/http-mpx/pkg/fetcher"
	"github.com/partyzanex/http-mpx/pkg/pool"
	"github.com/partyzanex/http-mpx/pkg/types"
)

// GetHandler returns api.Handler function.
func GetHandler(config Config, f fetcher.Interface) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		requests, count, err := parseRequest(&config, r)
		if err != nil {
			return err
		}

		var (
			res  = make(results, count)
			errs = make([]error, count)
			// creating context with cancel func for early exit from the pool
			poolCtx, poolCancel = context.WithCancel(r.Context())
			// creating a *pool.Pool instance for goroutines management
			wp = pool.New(config.Outgoing, pool.WithContext(poolCtx))
		)

		for i, request := range requests {
			t := task{
				Index:   i,
				Request: *request,
				Context: poolCtx,
				Cancel:  poolCancel,
				Errors:  errs,
				Results: res,
				Timeout: config.Timeout,
				Fetcher: f,
			}

			wp.Add(t.Worker)
		}

		wp.Wait()
		poolCancel() // release pool context

		if err := handleErrors(count, errs...); err != nil {
			return err
		}

		w.Header().Set("Content-Type", "application/json")

		return json.NewEncoder(w).Encode(res)
	}
}

// parseRequest returns []*types.Request from http.Request body and count.
func parseRequest(config *Config, r *http.Request) (requests []*types.Request, count int, err error) {
	// decode request
	err = json.NewDecoder(r.Body).Decode(&requests)
	if err != nil {
		return nil, 0, api.NewError(http.StatusBadRequest, MsgDecodeRequest).
			SetInternal(err)
	}

	count = len(requests)

	// check counts of requests
	if count > config.MaxURls {
		return nil, 0, api.NewError(http.StatusRequestEntityTooLarge, MsgLargeRequest)
	}

	return
}

// handleErrors join []errors to one error.
func handleErrors(count int, errs ...error) error {
	strErrors := make([]string, 0, count)

	for _, e := range errs {
		if e == nil {
			continue
		}

		strErrors = append(strErrors, e.Error())
	}

	if len(strErrors) > 0 {
		return api.NewError(http.StatusExpectationFailed, strings.Join(strErrors, "; "))
	}

	return nil
}
