package metrics

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	statusLabel = "status"
	methodLabel = "method"
	routeLabel  = "route"
)

var (
	requestDurationHistogram *prometheus.HistogramVec
	requestDurationSummary   *prometheus.SummaryVec
)

type Metrics interface {
	// Register provides a prometheus factory that can be employed to add custom metrics to the service
	Register(pf promauto.Factory)
}

// InitializeMetrics create a custom promauto instance with a local registry
// to avoid potential conflicts arising from packages using the default promauto registry
func InitializeMetrics(enableDefaultCollectors bool) (*prometheus.Registry, promauto.Factory) {
	reg := prometheus.NewRegistry()

	if enableDefaultCollectors {
		reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		reg.MustRegister(collectors.NewGoCollector())
	}

	return reg, promauto.With(reg)
}

// SetRequestMetrics register a set of metrics usefult to monitor the requests that are performed to the service
func setRequestMetrics(pf promauto.Factory) {
	requestDurationHistogram = pf.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "request duration in seconds",
			Buckets: []float64{0.05, 0.1, 0.5, 1, 3, 5, 10},
		},
		[]string{statusLabel, methodLabel, routeLabel},
	)
	requestDurationSummary = pf.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_summary_seconds",
			Help:       "request duration in seconds summary",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{statusLabel, methodLabel, routeLabel},
	)
}

func RequestStatus(pf promauto.Factory) func(http.Handler) http.Handler {
	setRequestMetrics(pf)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// avoid counting requests that does not belong to APIs
			if r.Method == http.MethodConnect || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			// default to status 200 to avoid empty values when WriteHeader
			// is not called to change the default status value 200 - OK
			httpResponse := httpResponseWriter{w, "200"}

			next.ServeHTTP(&httpResponse, r)

			end := time.Since(start).Seconds()
			// use path params patterns rather than actual value to avoid
			// generating too many different values for path label
			path := chi.RouteContext(r.Context()).RoutePattern()

			requestDurationHistogram.
				WithLabelValues(httpResponse.status, r.Method, path).
				Observe(end)
			requestDurationSummary.
				WithLabelValues(httpResponse.status, r.Method, path).
				Observe(end)
		})
	}
}
