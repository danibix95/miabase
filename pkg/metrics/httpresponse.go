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

func (hrw *httpResponseWriter) Header() http.Header {
	return hrw.writer.Header()
}

func (hrw *httpResponseWriter) Write(body []byte) (int, error) {
	return hrw.writer.Write(body)
}

func (hrw *httpResponseWriter) WriteHeader(statusCode int) {
	hrw.status = strconv.Itoa(statusCode)
	hrw.writer.WriteHeader(statusCode)
}

func (hrw *httpResponseWriter) Flush() {
	if f, ok := hrw.writer.(http.Flusher); ok {
		f.Flush()
	}
}
