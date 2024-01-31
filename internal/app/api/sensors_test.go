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

	// Create sensor using POST /sensors
	rr := createSensor(t, router, `
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
	`)
	// Should response w/201
	require.Equal(t, http.StatusCreated, rr.Code)

	// Check json response
	var res map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &res)
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

func TestGetSensorByName(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor using POST /sensors
	rr := createSensor(t, router, `
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
	`)
	// Should response w/201
	require.Equal(t, http.StatusCreated, rr.Code)

	// Get the sensor, using GET /sensors/:name
	handler := router.Handler()
	rr = httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sensors/abc123", nil)
	handler.ServeHTTP(rr, req)
	require.NoError(t, err)

	// Should respond with a 200
	require.Equal(t, http.StatusOK, rr.Code)

	// Unmarshal GET /sensors JSON response
	var res map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)

	// Should return the sensor that we created earlier
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

func TestGetSensorByNameMissing(t *testing.T) {
	router := NewSensorRouter()

	// Get a sensor that doesn't exist, using GET /sensors/:name
	handler := router.Handler()
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/sensors/not-a-sensor", nil)
	handler.ServeHTTP(rr, req)
	require.NoError(t, err)

	// Should respond with a 404
	require.Equal(t, http.StatusNotFound, rr.Code)

	// Unmarshal GET /sensors JSON response
	var res map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)

	// Should return the sensor that we created earlier
	require.Equal(t, map[string]interface{}{
		"error": "no sensor exists with name \"not-a-sensor\"",
	}, res)
}

func createSensor(t *testing.T, router *SensorRouter, sensorJSON string) *httptest.ResponseRecorder {
	handler := router.Handler()
	rr := httptest.NewRecorder()

	// Send POST /sensors request
	req, err := http.NewRequest("POST", "/sensors", bytes.NewBufferString(sensorJSON))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	return rr
}
