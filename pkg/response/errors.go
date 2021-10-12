package response

import (
	"net/http"

	"github.com/danibix95/zeropino"
	zpstd "github.com/danibix95/zeropino/middlewares/std"
)

type ErrorMessage struct {
	Msg  string `json:"message"`
	Code int    `json:"code,omitempty"`
}

func NotFound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	JSON(rw, ErrorMessage{Msg: "Route not found"})
}

func MethodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusMethodNotAllowed)
	JSON(rw, ErrorMessage{Msg: "Method not allowed"})
}

func InternalServerError(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	JSON(rw, ErrorMessage{Msg: "Generic server error"})
}

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
