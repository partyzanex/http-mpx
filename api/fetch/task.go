package fetch

import (
	"context"
	"encoding/json"
	"time"

	"github.com/partyzanex/http-mpx/pkg/fetcher"
	"github.com/partyzanex/http-mpx/pkg/types"
)

// results represents response body.
type results []json.RawMessage

// task represents a Worker params.
type task struct {
	// Index is an element index in slices Errors and Results
	Index int
	// Request is a copy of types.Request
	Request types.Request
	// Context is a parent context
	Context context.Context
	// Cancel is a cancel function for release a parent context
	Cancel context.CancelFunc
	// Errors is a slice for errors
	Errors []error
	// Results is a slice for results
	Results results
	// Timeout is max waiting time for Worker executing
	Timeout time.Duration
	// Fetcher instance of fetcher.Interface
	Fetcher fetcher.Interface
}

// Worker implements a pool.Worker.
func (t *task) Worker() {
	// create context with timeout via t.Context as a parent context
	ctx, cancel := context.WithTimeout(t.Context, t.Timeout)
	defer cancel()

	result, err := t.Fetcher.Fetch(ctx, t.Request)
	if err != nil {
		t.Errors[t.Index] = err
		t.Cancel() // cancel for parent context

		return
	}

	raw, err := json.Marshal(result)
	if err != nil {
		t.Errors[t.Index] = err
		t.Cancel() // cancel for parent context

		return
	}

	t.Results[t.Index] = raw
}
