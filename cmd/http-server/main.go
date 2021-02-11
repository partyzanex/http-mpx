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

	"github.com/partyzanex/http-mpx/internal/server"
	"github.com/partyzanex/http-mpx/pkg/fetcher/basehttp"
)

func main() {
	var (
		addr      = flag.String("addr", "0.0.0.0:3000", "address")
		outgoing  = flag.Int("outgoing", 4, "maximum of outgoing requests per connection")
		maxURLs   = flag.Int("max-urls", 20, "maximum of requested URLs per request")
		rateLimit = flag.Int64("rate-limit", 100, "rate limit")
		timeout   = flag.Duration("timeout", time.Second, "timeout for one request")
	)

	flag.Parse()

	srv := server.New(server.Config{
		Addr:      *addr,
		Outgoing:  *outgoing,
		MaxURls:   *maxURLs,
		RateLimit: *rateLimit,
		Timeout:   *timeout,
		Fetcher:   basehttp.New(nil),
	})

	// graceful shutdown
	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Println("Run server on", *addr)
	// run server
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("cannot run server: %s", err)
	}

	<-idleConnsClosed
}
