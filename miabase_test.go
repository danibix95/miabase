package miabase

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"

	"github.com/danibix95/miabase/pkg/response"
	"github.com/stretchr/testify/require"
)

func TestMiaBase(t *testing.T) {
	s := NewService()

	t.Run("Add route to plugin", func(t *testing.T) {
		// Add test handler
		message := map[string]interface{}{"msg": "welcome"}

		s.Plugin.Get("/greet", func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
			response.JSON(rw, message)
		})

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/greet", nil)
		response := executeRequest(t, req, s)

		require.Equal(t, http.StatusOK, response.Code, "Status codes mismatch")

		expectedResponse := message
		verifyJSONResponse(t, response, expectedResponse)
	})
}

// TestServiceStart verify that the bare bone service
// is able to start and to terminate gracefully
func TestServiceStart(t *testing.T) {
	s := NewService()

	go func() {
		time.Sleep(500 * time.Millisecond)
		s.SignalReceiver <- syscall.SIGTERM
	}()

	s.Start()
}

func executeRequest(t *testing.T, req *http.Request, s *Service) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	s.Plugin.ServeHTTP(rr, req)

	return rr
}

func verifyJSONResponse(t *testing.T, response *httptest.ResponseRecorder, expectedData interface{}) {
	var jsonData map[string]interface{}

	err := json.Unmarshal(response.Body.Bytes(), &jsonData)
	require.NoError(t, err)
	require.Equal(t, expectedData, jsonData)
}
