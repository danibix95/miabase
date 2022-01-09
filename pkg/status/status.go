package status

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Status interface {
	// Health returns an handler function that compute the liveness property of the service
	Health(name, version string) http.HandlerFunc
	// Ready returns an handler function that compute the readiness property of the service
	Ready(name, version string) http.HandlerFunc
	// CheckUp returns an handler function that verify the status of potential dependent services
	// and report
	CheckUp(name, version string) http.HandlerFunc
}

type Response struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
}

type DefaultStatus struct{}

func (ds DefaultStatus) Health(name, version string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		JSONResponse(rw, Response{Name: name, Version: version, Status: "OK"})
	}
}

func (ds DefaultStatus) Ready(name, version string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		JSONResponse(rw, Response{Name: name, Version: version, Status: "OK"})
	}
}

func (ds DefaultStatus) CheckUp(name, version string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		JSONResponse(rw, Response{Name: name, Version: version, Status: "OK"})
	}
}

func JSONResponse(w http.ResponseWriter, res Response) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := fmt.Fprintf(w, `{"name":%s,"version":%s,"status":KO}`, res.Name, res.Version); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
	}
}
