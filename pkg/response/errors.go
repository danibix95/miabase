package response

import (
	"net/http"
)

type ErrorMessage struct {
	Msg  string `json:"msg"`
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
