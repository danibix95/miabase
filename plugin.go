package miabase

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Plugin struct {
	Path   string
	router *chi.Mux
}

// NewPlugin create a new plugin that groups a set of routes under it
func NewPlugin(path string) *Plugin {
	p := new(Plugin)
	p.Path = path
	p.router = chi.NewRouter()

	return p
}

// AddRoute add a new endpoint to the plugin associated with the logic
// that should be executed when the route is called
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
		panic("selected HTTP method is not recognized")
	}
}

// Inject allow to test plugin routes by injecting the request and recording the response
func (p *Plugin) Inject(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
