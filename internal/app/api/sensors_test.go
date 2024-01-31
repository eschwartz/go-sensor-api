package api

import (
	"bytes"
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

func TestCreateSensor(t *testing.T) {
	router := NewSensorRouter()
	handler := router.Handler()
	rr := httptest.NewRecorder()

	// Prepare Sensor JSON body
	sensorJSON := bytes.NewBuffer([]byte(`
		{
		  "name": "abc123",
		  "lat": 44.916241209323736,
		  "lon": -93.21112681214602,
		  "tags": [
			"x",
			"y",
			"z"
		  ]
		}
	`))

	// Send POST /sensors request
	req, err := http.NewRequest("POST", "/sensors", sensorJSON)
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	// Should response w/201
	require.Equal(t, http.StatusCreated, rr.Code)

	// Check json response
	var res map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)

	require.Equal(t, map[string]interface{}{
		"data": map[string]interface{}{
			"name": "abc123",
			"lat":  44.916241209323736,
			"lon":  -93.21112681214602,
			"tags": []interface{}{
				"x",
				"y",
				"z",
			},
		},
	}, res)
}
