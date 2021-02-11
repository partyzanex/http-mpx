package server

import (
	"time"

	"github.com/partyzanex/http-mpx/pkg/fetcher"
)

// Config represent a configuration for http.Server.
type Config struct {
	// Addr is a listen address
	Addr string
	// Outgoing is a count of maximum outgoing requests
	// for each incoming request
	Outgoing int
	// MaxURls is a count of maximum URLs
	// in each incoming request
	MaxURls int
	// RateLimit is a count of maximum incoming
	RateLimit int
	Timeout   time.Duration

	Fetcher fetcher.Interface
}
