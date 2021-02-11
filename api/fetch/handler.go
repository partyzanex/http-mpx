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

func NewHandler(config Config, f fetcher.Interface) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var requests types.Requests

		err := json.NewDecoder(r.Body).Decode(&requests)
		if err != nil {
			return api.NewError(http.StatusBadRequest, "cannot decode request").
				SetInternal(err)
		}

		defer r.Body.Close()

		if len(requests) > config.MaxURls {
			return api.NewError(http.StatusRequestEntityTooLarge, "large count of requests")
		}

		var (
			p       = pool.New(config.Outgoing)
			count   = len(requests)
			results = make([]*types.Result, count)
			errs    = make([]error, count)
		)

		for i, request := range requests {
			n := i
			req := *request

			p.Add(func() {
				ctx, cancel := context.WithTimeout(r.Context(), config.Timeout)
				defer cancel()

				result, err := f.Fetch(ctx, req)
				if err != nil {
					errs[n] = err
					return
				}

				results[n] = result
			})
		}

		p.Wait()

		strErrors := make([]string, 0)

		for _, e := range errs {
			if e == nil {
				continue
			}

			strErrors = append(strErrors, e.Error())
		}

		if len(strErrors) > 0 {
			return api.NewError(http.StatusExpectationFailed, strings.Join(strErrors, "; "))
		}

		return json.NewEncoder(w).Encode(results)
	}
}
