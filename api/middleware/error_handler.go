package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/partyzanex/http-mpx/api"
)

const errFormat = "error: code=%d message=%s internal=%v"

// ErrorHandler implements the api.Middleware function
// that handles the error returned by the wrapped handler
// and write it to response in JSON format.
func ErrorHandler(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := next(w, r)
		if err != nil {
			// handle error
			var e *api.Error

			switch er := err.(type) {
			case *api.Error:
				e = er
			default:
				e = &api.Error{
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				}
			}

			// write error to log
			defer log.Printf(errFormat, e.Code, e.Message, e.Internal)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(e.Code)

			return json.NewEncoder(w).Encode(e)
		}

		return err
	}
}
