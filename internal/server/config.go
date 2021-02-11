package server

import (
	"time"

	"github.com/partyzanex/http-mpx/pkg/fetcher"
)

type Config struct {
	Addr string

	Outgoing  int
	MaxURls   int
	RateLimit int64
	Timeout   time.Duration

	Fetcher fetcher.Interface
}
