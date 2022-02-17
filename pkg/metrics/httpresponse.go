package metrics

import (
	"net/http"
	"strconv"
)

// httpResponseWriter uses extension interface pattern to decorate the status code to the response
type httpResponseWriter struct {
	writer http.ResponseWriter
	status string
}

// Header return the Header map of the wrapped http ResponseWriter
func (hrw *httpResponseWriter) Header() http.Header {
	return hrw.writer.Header()
}

// Write execute the Write method on the wrapped http ResponseWriter
func (hrw *httpResponseWriter) Write(body []byte) (int, error) {
	return hrw.writer.Write(body)
}

// WriteHeader store the statusCode within the response wrapper and then
// execute the http ResponseWriter's WriteHeade method
func (hrw *httpResponseWriter) WriteHeader(statusCode int) {
	hrw.status = strconv.Itoa(statusCode)
	hrw.writer.WriteHeader(statusCode)
}

// Flush when the wrapped http ResponseWriter implements the Flusher interface,
// it execute the Flush method
func (hrw *httpResponseWriter) Flush() {
	if f, ok := hrw.writer.(http.Flusher); ok {
		f.Flush()
	}
}
