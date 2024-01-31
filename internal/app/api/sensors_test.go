package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	router := NewSensorRouter()
	handler := router.Handler()
	rr := httptest.NewRecorder()

	// Send GET /health request
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	// Should response w/200
	require.Equal(t, http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	var res map[string]bool
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)

	require.Equal(t, true, res["ok"])
}
