package fetcher

import (
	"context"

	"github.com/partyzanex/http-mpx/pkg/types"
)

type Interface interface {
	Fetch(ctx context.Context, req types.Request) (*types.Result, error)
}
