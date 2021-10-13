package miabase

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danibix95/miabase/pkg/handlers"
	"github.com/danibix95/miabase/pkg/response"
	"github.com/danibix95/zeropino"
	zpstd "github.com/danibix95/zeropino/middlewares/std"
	"github.com/go-chi/chi/v5"
	"github.com/mia-platform/configlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type Service struct {
	Name           string
	Version        string
	router         *chi.Mux
	Plugin         *chi.Mux
	Logger         *zerolog.Logger
	HealthHandler  http.HandlerFunc
	ReadyHandler   http.HandlerFunc
	CheckupHandler http.HandlerFunc
	signalReceiver chan os.Signal
}

func LoadEnv(c []configlib.EnvConfig, env interface{}) {
	if err := configlib.GetEnvVariables(c, &env); err != nil {
		panic(err.Error())
	}
}

func NewService(name, version, logLevel string) *Service {
	s := new(Service)
	s.router = chi.NewRouter()
	s.Plugin = chi.NewRouter()
	s.Name = name
	s.Version = version

	logger, err := zeropino.Init(zeropino.InitOptions{Level: logLevel})
	if err != nil {
		panic(err.Error())
	}
	s.Logger = logger

	s.signalReceiver = make(chan os.Signal, 1)

	return s
}

// Start launch the configured service,
// mounting customized plugin and starting the webserver
func (s *Service) Start(httpPort int) {
	s.addErrorsHandlers()

	s.addStatusRoutes()

	s.router.Group(func(r chi.Router) {
		r.Use(zpstd.RequestLogger(s.Logger, []string{"/-/"}))

		r.Mount("/", s.Plugin)
	})

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", httpPort), Handler: s.router}

	runWithGracefulShutdown(server, s.Logger, s.signalReceiver)
}

func (s *Service) Stop() {
	s.signalReceiver <- syscall.SIGTERM
}

func (s *Service) addErrorsHandlers() {
	s.router.Use(response.PanicManager)
	s.router.NotFound(response.NotFound)
	s.router.MethodNotAllowed(response.MethodNotAllowed)
}

func (s *Service) addStatusRoutes() {
	s.router.Group(func(r chi.Router) {
		statusAndMetricsRouter := chi.NewRouter()

		if s.HealthHandler != nil {
			statusAndMetricsRouter.Get("/healthz", s.HealthHandler)
		} else {
			statusAndMetricsRouter.Get("/healthz", handlers.Health(s.Name, s.Version))
		}

		if s.ReadyHandler != nil {
			statusAndMetricsRouter.Get("/ready", s.ReadyHandler)
		} else {
			statusAndMetricsRouter.Get("/ready", handlers.Ready(s.Name, s.Version))
		}

		if s.CheckupHandler != nil {
			statusAndMetricsRouter.Get("/check-up", s.CheckupHandler)
		} else {
			statusAndMetricsRouter.Get("/check-up", handlers.CheckUp(s.Name, s.Version))
		}

		statusAndMetricsRouter.Handle("/metrics", promhttp.Handler())

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
