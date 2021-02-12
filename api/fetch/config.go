package fetch

import (
	"time"
)

// Config represent a configuration for http.Server.
type Config struct {
	// Outgoing is a count of maximum outgoing requests
	// for each incoming request
	Outgoing int
	// MaxURls is a count of maximum URLs
	// in each incoming request
	MaxURls int
	// Timeout is a wait time for each outgoing request
	Timeout time.Duration
}

const (
	MsgDecodeRequest = "cannot decode request"
	MsgLargeRequest  = "large count of requests"
)
