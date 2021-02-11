package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/partyzanex/http-mpx/pkg/limiter"
	"github.com/partyzanex/http-mpx/pkg/pool"
	"github.com/partyzanex/http-mpx/pkg/types"
)

const errFormat = "error: code=%d message=%s internal=%v"

// Handler implements of http.Handler interface.
type Handler struct {
	// Config is embed Config
	Config
	// limiter for calls count restriction
	limiter limiter.Limiter
}

// ServeHTTP implements of http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// recover
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("recovered", rec)
		}
	}()

	// apply limiter
	if !h.limiter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	// take a limit of calls
	release := h.limiter.Take()
	defer release()

	// check HTTP-method
	if r.Method != http.MethodPost {
		h.writeError(w, &Error{
			Code:    http.StatusMethodNotAllowed,
			Message: "method not allowed",
		})
		return
	}

	// write Content-Type header
	w.Header().Add("Content-Type", "application/json")

	// execute request
	err := h.Handle(w, r)
	if err != nil {
		// handle error
		var logErr *Error

		switch er := err.(type) {
		case *Error:
			logErr = er
		default:
			logErr = &Error{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		}

		h.writeError(w, logErr)
		return
	}
}

// writeError writes error in JSON to response
func (*Handler) writeError(w http.ResponseWriter, e *Error) {
	defer log.Printf(errFormat, e.Code, e.Message, e.Internal)

	w.WriteHeader(e.Code)

	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Println("cannot encode to json:", err)
	}
}

// Handle is handler of requests
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	var requests types.Requests

	err := json.NewDecoder(r.Body).Decode(&requests)
	if err != nil {
		return NewError(http.StatusBadRequest, "cannot decode request").
			SetInternal(err)
	}

	defer r.Body.Close()

	if len(requests) > h.MaxURls {
		return NewError(http.StatusRequestEntityTooLarge, "large count of requests")
	}

	var (
		p       = pool.New(h.Outgoing)
		count   = len(requests)
		results = make([]*types.Result, count)
		errs    = make([]error, count)
	)

	for i, request := range requests {
		n := i
		req := *request

		p.Add(func() {
			ctx, cancel := context.WithTimeout(r.Context(), h.Timeout)
			defer cancel()

			result, err := h.Fetcher.Fetch(ctx, req)
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
		return NewError(http.StatusExpectationFailed, strings.Join(strErrors, "; "))
	}

	return json.NewEncoder(w).Encode(results)
}
