package miabase

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/danibix95/miabase/pkg/response"
	"github.com/stretchr/testify/require"
)

func TestMiaBase(t *testing.T) {
	t.Run("Add route to plugin", func(t *testing.T) {
		s := NewService()
		// Add test handler
		message := map[string]interface{}{"message": "welcome"}

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

// TestPanicHandler verifies that a service
// is able to handle panics returning Internal Server Error
func TestPanicHandler(t *testing.T) {
	t.Run("Handle panic correctly", func(t *testing.T) {
		s := NewService()
		s.Plugin.Get("/panic", func(rw http.ResponseWriter, r *http.Request) {
			panic("it should not die")
		})

		go func() {
			time.Sleep(300 * time.Millisecond)
			s.Stop()
		}()

		s.Start()

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/panic", nil)

		response := httptest.NewRecorder()
		s.router.ServeHTTP(response, req)

		require.Equal(t, http.StatusInternalServerError, response.Code, "Status codes mismatch")
	})
}

// TestServiceStart verifies that the bare bone service
// is able to start and to terminate gracefully
func TestServiceStart(t *testing.T) {
	s := NewService()

	go func() {
		time.Sleep(300 * time.Millisecond)
		s.Stop()
	}()

	s.Start()
}

func executeRequest(t *testing.T, req *http.Request, s *Service) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	s.router.Mount("/", s.Plugin)
	s.router.ServeHTTP(rr, req)

	return rr
}

func verifyJSONResponse(t *testing.T, response *httptest.ResponseRecorder, expectedData interface{}) {
	var jsonData map[string]interface{}

	err := json.Unmarshal(response.Body.Bytes(), &jsonData)
	require.NoError(t, err)
	require.Equal(t, expectedData, jsonData)
}
