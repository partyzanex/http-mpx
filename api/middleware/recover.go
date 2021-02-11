package middleware

import (
	"fmt"
	"net/http"

	"github.com/partyzanex/http-mpx/api"
)

const recoverFormat = "recovered: %v"

// Recover middleware.
func Recover(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		// recover
		defer func() {
			if rec := recover(); rec != nil {
				err = api.NewError(http.StatusInternalServerError, "").
					SetInternal(fmt.Errorf(recoverFormat, rec))
			}
		}()

		err = next(w, r)

		return
	}
}
