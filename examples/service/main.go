package main

import (
	"fmt"
	"net/http"

	"github.com/danibix95/miabase"
	"github.com/danibix95/miabase/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/mia-platform/configlib"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// ----- Environment Section -----
type Environment struct {
	LogLevel string
	HTTPPort int
}

var envVariablesConfig = []configlib.EnvConfig{
	{
		Key:          "LOG_LEVEL",
		Variable:     "LogLevel",
		DefaultValue: "info",
	},
	{
		Key:          "HTTP_PORT",
		Variable:     "HTTPPort",
		DefaultValue: "3000",
	},
}

// ------- Metrics Section -------
var (
	greetinsCounter prometheus.Counter
)

type customMetrics struct{}

func (cm customMetrics) Register(pf promauto.Factory) {
	greetinsCounter = pf.NewCounter(prometheus.CounterOpts{
		Name: "greetings_total",
		Help: "count the number of times that the service greets a user",
	})
}

// -------------------------------

func main() {
	var env Environment
	miabase.LoadEnv(envVariablesConfig, &env)

	service := miabase.NewService(miabase.ServiceOpts{
		Name:           "example",
		Version:        "v0.0.1",
		LogLevel:       env.LogLevel,
		StatusManager:  nil,
		MetricsManager: customMetrics{},
	})

	plugin := miabase.NewPlugin("/")
	plugin.AddRoute("GET", "/greet", func(rw http.ResponseWriter, r *http.Request) {
		// use the custom metric
		greetinsCounter.Inc()

		response.JSON(rw, map[string]string{"message": "welcome"})
	})

	plugin.AddRoute("GET", "/ciaone/{who}", func(rw http.ResponseWriter, r *http.Request) {
		who := chi.URLParam(r, "who")

		response.JSON(rw, map[string]string{"message": fmt.Sprintf("ciaone %s", who)})
	})

	service.Register(plugin)

	service.Start(env.HTTPPort)
}
