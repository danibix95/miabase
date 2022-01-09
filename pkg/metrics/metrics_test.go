package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/stretchr/testify/require"
)

type metricsTest struct{}

func (mt metricsTest) Register(pf promauto.Factory) {
	pf.NewCounter(prometheus.CounterOpts{Name: "counter_test_total"})
}

func TestMetrics(t *testing.T) {
	t.Run("verify custom metrics manager implements Metrics interface", func(t *testing.T) {
		require.NotPanics(t, func() {
			var _ Metrics = new(metricsTest)
		})
	})
}

func TestInitializeMetrics(t *testing.T) {
	t.Run("create metrics factory without any default collector", func(t *testing.T) {
		metricsRegistry, metricsFactory := InitializeMetrics(false)
		require.NotNil(t, metricsRegistry)
		require.NotPanics(t, func() {
			metricsFactory.NewCounter(prometheus.CounterOpts{Name: "process_cpu_seconds_total"})
		}, "metric does not exists, so that it is fine to create it")
	})

	t.Run("create metrics factory with default collectors (GO and Processes)", func(t *testing.T) {
		metricsRegistry, metricsFactory := InitializeMetrics(true)
		require.NotNil(t, metricsRegistry)
		require.Panics(t, func() {
			metricsFactory.NewCounter(prometheus.CounterOpts{Name: "process_cpu_seconds_total"})
		}, "metric already exists, so that creating a new one cause the code to panic")
	})
}

func TestRequestMetrics(t *testing.T) {
	t.Run("verify default metrics are registered", func(t *testing.T) {
		// check that pointers are nil before execution
		require.Nil(t, requestDurationHistogram)
		require.Nil(t, requestDurationSummary)

		reg := prometheus.NewPedanticRegistry()
		promFactory := promauto.With(reg)

		require.NotPanics(t, func() {
			SetRequestMetrics(promFactory)
		}, "metrics are registered correctly")

		// an object has been assigned to the pointers
		require.NotNil(t, requestDurationHistogram)
		require.NotNil(t, requestDurationSummary)

		require.NotPanics(t, func() {
			requestDurationHistogram.WithLabelValues("200", "GET", "/greetings").Observe(0.07)
			requestDurationSummary.WithLabelValues("200", "GET", "/greetings").Observe(0.07)
		}, "metrics can be employed to observe some values")
	})
}

func TestRequestMetricsMiddleware(t *testing.T) {
	t.Run("", func(t *testing.T) {

	})
}
