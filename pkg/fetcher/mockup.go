package fetcher

import (
	"context"
	"errors"

	"github.com/partyzanex/http-mpx/pkg/types"
)

// Mockup implements Interface.
type Mockup struct {
	FetchMethod func(ctx context.Context, req types.Request) (*types.Result, error)
}

func (m *Mockup) Fetch(ctx context.Context, req types.Request) (*types.Result, error) {
	if m.FetchMethod == nil {
		return nil, errors.New("not implemented")
	}

	return m.FetchMethod(ctx, req)
}
