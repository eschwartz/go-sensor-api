package api

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	router := NewSensorRouter()

	rr := httptest.NewRecorder()

	handler := router.Handler()
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	var res map[string]bool
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)

	require.Equal(t, true, res["ok"])
}
