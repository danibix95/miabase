package status

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	name    = "service-name"
	version = "0.0.1"
)

func TestDefaultStatus(t *testing.T) {
	ds := &DefaultStatus{}
	expectedResponse := fmt.Sprintf(`{"name":"%s","version":"%s","status":"OK"}`, name, version)

	t.Run("Health handler", func(t *testing.T) {
		healthHandler := ds.Health(name, version)

		verifyStatusRequest(t, healthHandler, "/-/healthz", expectedResponse)
	})

	t.Run("Ready handler", func(t *testing.T) {
		healthHandler := ds.Ready(name, version)

		verifyStatusRequest(t, healthHandler, "/-/ready", expectedResponse)
	})

	t.Run("CheckUp handler", func(t *testing.T) {
		healthHandler := ds.CheckUp(name, version)

		verifyStatusRequest(t, healthHandler, "/-/check-up", expectedResponse)
	})
}

func verifyStatusRequest(t *testing.T, h http.HandlerFunc, endpoint, expectedResponse string) {
	t.Helper()

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, expectedResponse, strings.TrimSpace(rr.Body.String()))
}
