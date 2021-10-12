package miabase

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danibix95/miabase/pkg/response"
	"github.com/danibix95/zeropino"
	zpstd "github.com/danibix95/zeropino/middlewares/std"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	router         *chi.Mux
	Plugin         *chi.Mux
	Logger         *zerolog.Logger
	signalReceiver chan os.Signal
}

func NewService() *Service {
	s := new(Service)
	s.router = chi.NewRouter()
	s.Plugin = chi.NewRouter()

	logger, err := zeropino.Init(zeropino.InitOptions{Level: "info"})
	if err != nil {
		panic(err.Error())
	}
	s.Logger = logger

	s.signalReceiver = make(chan os.Signal, 1)

	return s
}

// Start launch the configured service,
// mounting customized plugin and starting the webserver
func (s *Service) Start() {
	s.addErrorsHandlers()

	s.router.Group(func(r chi.Router) {
		r.Use(zpstd.RequestLogger(s.Logger, []string{"/-/"}))

		r.Mount("/", s.Plugin)
	})

	server := &http.Server{Addr: "0.0.0.0:3000", Handler: s.router}

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
	log.Info().Msg("server listening at http://localhost:3000")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server closed unexpectedly")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
