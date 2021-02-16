package api

import (
	"log"
	"net/http"
)

// Handler implements the http.Handler interface.
type Handler func(w http.ResponseWriter, r *http.Request) error

// ServeHTTP implements the http.Handler interface.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		log.Println(err)
	}
}

// Middleware function for wrapping Handler.
type Middleware func(next Handler) Handler

// Wrap wraps Handler with Middleware functions.
func Wrap(handler Handler, middlewares ...Middleware) Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
