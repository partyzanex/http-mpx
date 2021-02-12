package fetcher

import (
	"context"

	"github.com/partyzanex/http-mpx/pkg/types"
)

// Interface is a fetcher interface.
type Interface interface {
	// Fetch should be takes a types.Request, execute HTTP-request
	// and returns *types.Result or error
	Fetch(ctx context.Context, req types.Request) (*types.Result, error)
}
