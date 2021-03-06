package miabase

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danibix95/miabase/pkg/metrics"
	"github.com/danibix95/miabase/pkg/response"
	"github.com/danibix95/miabase/pkg/status"
	"github.com/danibix95/zeropino"
	zpstd "github.com/danibix95/zeropino/middlewares/std"
	"github.com/go-chi/chi/v5"
	"github.com/mia-platform/configlib"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

// Service is the main structure that contains all the service details,
// the methods to attach custom plugins and the ones to start it
type Service struct {
	name            string
	version         string
	router          *chi.Mux
	plugins         []*Plugin
	statusManager   status.Status
	signalReceiver  chan os.Signal
	metricsRegistry *prometheus.Registry
	metricsFactory  promauto.Factory
	// Logger a zerolog instance that can be employed to log service details within plugins
	Logger *zerolog.Logger
}

// ServiceOpts defines which options are needed to customize a Service initialization
type ServiceOpts struct {
	// Name represents the designation emplyoyed to indentify the service's deploy
	Name string
	// Version is a semver string that represents the current version of deployed service
	Version string
	// LogLevel is a string indicating the minimum log level that is shown on the standard out
	LogLevel string
	// StatusManager is an interface providing the three status routes handlers
	StatusManager status.Status
	// MetricsManager is an interface providing a method to register custom metrics in the service registry
	MetricsManager metrics.Metrics
}

func LoadEnv(c []configlib.EnvConfig, env interface{}) {
	if err := configlib.GetEnvVariables(c, &env); err != nil {
		panic(err.Error())
	}
}

// NewService instantiate a Service which can be employed to connect custom plugin
// and start listening on defined endpoints
func NewService(opts ServiceOpts) *Service {
	s := new(Service)
	s.router = chi.NewRouter()
	s.name = opts.Name
	s.version = opts.Version

	logger, err := zeropino.Init(zeropino.InitOptions{Level: opts.LogLevel})
	if err != nil {
		panic(err.Error())
	}
	s.Logger = logger

	if opts.StatusManager == nil {
		s.statusManager = status.DefaultStatus{}
	} else {
		s.statusManager = opts.StatusManager
	}

	s.metricsRegistry, s.metricsFactory = metrics.InitializeMetrics(true)
	if opts.MetricsManager != nil {
		opts.MetricsManager.Register(s.metricsFactory)
	}

	s.signalReceiver = make(chan os.Signal, 1)

	return s
}

// Register include the new plugin into the set of plugins that the service must load.
func (s *Service) Register(plugin *Plugin) {
	s.plugins = append(s.plugins, plugin)
}

// Start launch the configured service,
// mounting customized plugin and starting the webserver
func (s *Service) Start(httpPort int) {
	s.setupServicePlugins()

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", httpPort), Handler: s.router}

	runWithGracefulShutdown(server, s.Logger, s.signalReceiver)
}

// Stop terminates service webserver execution
func (s *Service) Stop() {
	s.signalReceiver <- syscall.SIGTERM
}

// Inject allows to test an endpoint by passing a request and response recorder
func (s *Service) Inject(w http.ResponseWriter, r *http.Request) {
	s.setupServicePlugins()

	s.router.ServeHTTP(w, r)
}

func (s *Service) setupServicePlugins() {
	s.addErrorsHandlers()
	s.router.Use(metrics.RequestStatus(s.metricsFactory))
	s.addStatusRoutes()

	s.router.Group(func(r chi.Router) {
		r.Use(zpstd.RequestLogger(s.Logger, []string{"/-/"}))

		for _, plugin := range s.plugins {
			r.Mount(plugin.Path, plugin.router)
		}
	})
}

func (s *Service) addErrorsHandlers() {
	s.router.Use(response.PanicManager)
	s.router.NotFound(response.NotFound)
	s.router.MethodNotAllowed(response.MethodNotAllowed)
}

func (s *Service) addStatusRoutes() {
	s.router.Group(func(r chi.Router) {
		statusAndMetricsRouter := chi.NewRouter()

		statusAndMetricsRouter.Get("/healthz", s.statusManager.Health(s.name, s.version))
		statusAndMetricsRouter.Get("/ready", s.statusManager.Ready(s.name, s.version))
		statusAndMetricsRouter.Get("/check-up", s.statusManager.CheckUp(s.name, s.version))

		statusAndMetricsRouter.Handle("/metrics", promhttp.HandlerFor(s.metricsRegistry, promhttp.HandlerOpts{}))

		r.Mount("/-/", statusAndMetricsRouter)
	})
}

func runWithGracefulShutdown(srv *http.Server, log *zerolog.Logger, sig chan os.Signal) {
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Error().Msg("graceful shutdown timed out.. forcing exit")
			}
		}()

		// Trigger graceful shutdown
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err).Msg("server shutdown did not work as expected")
		}

		serverStopCtx()
		shutdownStopCtx()
	}()

	// Run the server
	log.Info().Msg(fmt.Sprintf("server listening at %s", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server closed unexpectedly")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
