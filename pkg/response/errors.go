package response

import (
	"net/http"

	"github.com/danibix95/zeropino"
	zpstd "github.com/danibix95/zeropino/middlewares/std"
)

type errorMessage struct {
	Msg  string `json:"message"`
	Code int    `json:"code,omitempty"`
}

// NotFound is an http handler that return a JSON response
// when requested resource is not found at the current route
func NotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	JSON(rw, errorMessage{Msg: "Route not found"})
}

// MethodNotAllowed is an http handler that return a JSON response
// when the method of current request has not been defined for current route
func MethodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
	JSON(rw, errorMessage{Msg: "Method not allowed"})
}

// InternalServerError is an http handler that returns a JSON response
// when an error that can not be managed by the handler is encountered during requests handling
func InternalServerError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	JSON(rw, errorMessage{Msg: "Generic server error"})
}

// PanicManager return a middleware function that recover service from
// panic situations, by returning an Interal Server Error response
func PanicManager(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				logEntry := zpstd.Get(r.Context())

				if logEntry != nil {
					logEntry.Error().Stack().Msg("recovered from panic")
				} else {
					zeropino.InitDefault().Error().Stack().Msg("recovered from panic")
				}

				InternalServerError(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
