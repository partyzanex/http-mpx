package middleware

import (
	"net/http"

	"github.com/partyzanex/http-mpx/api"
)

// AllowedMethods creates api.Middleware function,
// that returns api.Handler with applied filter by HTTP-method.
func AllowedMethods(methods ...string) api.Middleware {
	return func(next api.Handler) api.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			// filtering HTTP-method
			for _, method := range methods {
				if method == r.Method {
					return next(w, r)
				}
			}

			// return *api.Error
			return api.NewError(http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
