package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type StatusResponse struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
}

func Health(name, version string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		JSONResponse(rw, StatusResponse{Name: name, Version: version, Status: "OK"})
	}
}

func Ready(name, version string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		JSONResponse(rw, StatusResponse{Name: name, Version: version, Status: "OK"})
	}
}

func CheckUp(name, version string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		JSONResponse(rw, StatusResponse{Name: name, Version: version, Status: "OK"})
	}
}

func JSONResponse(w http.ResponseWriter, res StatusResponse) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err := fmt.Fprintf(w, `{"name":%s,"version":%s,"status":KO}`, res.Name, res.Version); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
	}
}
