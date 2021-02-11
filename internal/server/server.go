package server

import (
	"net/http"

	"github.com/partyzanex/http-mpx/pkg/limiter"
)

// New creates a new *http.Server instance
// with applied configuration (Config).
func New(config Config) *http.Server {
	return &http.Server{
		Addr: config.Addr,
		Handler: &Handler{
			Config:  config,
			limiter: limiter.New(config.RateLimit),
		},
	}
}
