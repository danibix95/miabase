package miabase

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	router *chi.Mux
	Plugin *chi.Mux
}

func NewService() *Service {
	s := new(Service)
	s.router = chi.NewRouter()
	s.Plugin = chi.NewRouter()

	return s
}

func (s *Service) Start() {
	s.router.Group(func(r chi.Router) {
		r.Mount("/", s.Plugin)
	})

	server := &http.Server{Addr: "0.0.0.0:3000", Handler: s.router}
	_ = server.ListenAndServe()
}
