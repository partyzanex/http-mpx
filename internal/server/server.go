package server

import (
	"net/http"
	"time"

	"github.com/partyzanex/http-mpx/pkg/limiter"
)

func New(config Config) *http.Server {
	return &http.Server{
		Addr: config.Addr,
		Handler: &Handler{
			Config:  config,
			limiter: limiter.New(time.Second, config.RateLimit),
		},
	}
}
