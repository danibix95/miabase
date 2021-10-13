package miabase

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Plugin struct {
	Path   string
	router *chi.Mux
}

func NewPlugin(path string) *Plugin {
	p := new(Plugin)
	p.Path = path
	p.router = chi.NewRouter()

	return p
}

func (p *Plugin) AddRoute(method, path string, handler http.HandlerFunc) {
	switch method {
	case "GET":
		p.router.Get(path, handler)
	case "POST":
		p.router.Post(path, handler)
	case "PUT":
		p.router.Put(path, handler)
	case "PATCH":
		p.router.Patch(path, handler)
	case "DELETE":
		p.router.Delete(path, handler)
	default:
		panic("selected method it is not recognized")
	}
}
