package miabase

import (
	"net/http"

	"github.com/danibix95/zeropino"
	zpstd "github.com/danibix95/zeropino/middlewares/std"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	router *chi.Mux
	Plugin *chi.Mux
	Logger *zerolog.Logger
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

	return s
}

func (s *Service) Start() {
	s.router.Group(func(r chi.Router) {
		r.Use(zpstd.RequestLogger(s.Logger, []string{"/-/"}))

		r.Mount("/", s.Plugin)
	})

	server := &http.Server{Addr: "0.0.0.0:3000", Handler: s.router}
	s.Logger.Info().Msg("Server listening on port :3000")
	_ = server.ListenAndServe()
}
