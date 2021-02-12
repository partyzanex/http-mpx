package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/partyzanex/http-mpx/api"
	"github.com/partyzanex/http-mpx/api/fetch"
	"github.com/partyzanex/http-mpx/api/middleware"
	"github.com/partyzanex/http-mpx/pkg/fetcher/basehttp"
)

func main() {
	var (
		addr      = flag.String("addr", "0.0.0.0:3000", "address")
		outgoing  = flag.Int("outgoing", 4, "maximum of outgoing requests per connection")
		maxURLs   = flag.Int("max-urls", 20, "maximum of requested URLs per request")
		rateLimit = flag.Int("rate-limit", 100, "rate limit")
		timeout   = flag.Duration("timeout", time.Second, "timeout for one request")
	)

	flag.Parse()

	fetchConfig := fetch.Config{
		Outgoing: *outgoing,
		MaxURls:  *maxURLs,
		Timeout:  *timeout,
	}
	fetchHandler := api.Wrap(
		fetch.GetHandler(fetchConfig, basehttp.New(nil)),
		middleware.Logger, middleware.ConcurrentLimiter(*rateLimit),
		middleware.AllowedMethods(http.MethodPost),
		middleware.Recover, middleware.ErrorHandler,
	)

	server := http.Server{
		Addr:    *addr,
		Handler: fetchHandler,
	}

	// graceful shutdown
	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Println("Listen on", *addr)
	// run server
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("cannot run server: %s", err)
	}

	<-idleConnsClosed
}
