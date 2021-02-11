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

type Handler struct {
	Config

	limiter limiter.Limiter
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println("recovered", rec)
		}
	}()

	if !h.limiter.Allow() {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		e := json.NewEncoder(w).Encode(Error{
			Message: "method not allowed",
		})
		if e != nil {
			log.Println("cannot encode to json:", e)
		}

		return
	}

	w.Header().Add("Content-Type", "application/json")

	err := h.Handle(w, r)
	if err != nil {
		var logErr *Error

		defer func() {
			log.Printf("error: code=%d message=%s internal=%v", logErr.Code, logErr.Message, logErr.Internal)
		}()

		switch er := err.(type) {
		case *Error:
			logErr = er
		default:
			logErr = &Error{
				Message: err.Error(),
				Code:    http.StatusInternalServerError,
			}
		}

		w.WriteHeader(logErr.Code)

		errWr := json.NewEncoder(w).Encode(logErr)
		if errWr != nil {
			log.Println("cannot encode to json:", errWr)
		}

		return
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) error {
	var requests types.Requests

	err := json.NewDecoder(r.Body).Decode(&requests)
	if err != nil {
		return NewError(http.StatusBadRequest, "cannot decode request").
			SetInternal(err)
	}

	if len(requests) > h.MaxURls {
		return NewError(http.StatusRequestEntityTooLarge, "large count of requests")
	}

	p := pool.New(h.Outgoing)
	count := len(requests)
	results := make([]*types.Result, count)
	errs := make([]error, count)

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
