package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io"
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

func TestCreateSensor_Invalid(t *testing.T) {
	// TODO
	require.True(t, false, "TODO")
}

func TestCreateSensor_StoreFailure(t *testing.T) {
	// TODO
	// Maybe make a store that fails on every method?
	require.True(t, false, "TODO")
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

func TestGetSensorByName_Missing(t *testing.T) {
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

	// Should return an error
	require.Equal(t, map[string]interface{}{
		"error": "no sensor resource exists: not-a-sensor",
	}, res)
}

func TestUpdateSensorByName(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor to work with using POST /sensors
	rr := httpRequest(t, router, "POST", "/sensors", `
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

	// Update the sensor, using PUT /sensors/abc123
	rr = httpRequest(t, router, "PUT", "/sensors/abc123", `
		{
		  "name": "abc123",
		  "lat": -36.8779565276809,
		  "lon": 174.7881226266269744,
		  "tags": [
			"a",
			"b",
			"c"
		  ]
		}
	`)
	require.Equal(t, http.StatusOK, rr.Code)

	// Inspect the PUT response
	putRes := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"data": map[string]interface{}{
			"name": "abc123",
			"lat":  -36.8779565276809,
			"lon":  174.7881226266269744,
			"tags": []interface{}{
				"a",
				"b",
				"c",
			},
		},
	}, putRes, "PUT response")

	// Retrieve the updated sensor, using GET /sensors/abc123
	rr = httpRequest(t, router, "GET", "/sensors/abc123", "")
	require.Equal(t, http.StatusOK, rr.Code)
	getRes := unmarshalResponseJSON(t, rr)

	require.Equal(t, map[string]interface{}{
		"data": map[string]interface{}{
			"name": "abc123",
			"lat":  -36.8779565276809,
			"lon":  174.7881226266269744,
			"tags": []interface{}{
				"a",
				"b",
				"c",
			},
		},
	}, getRes, "GET response")
}

func TestUpdateSensorByName_Missing(t *testing.T) {
	router := NewSensorRouter()

	// Update a sensor that doesn't exist, using PUT /sensors/not-a-sensor
	rr := httpRequest(t, router, "PUT", "/sensors/not-a-sensor", `
		{
		  "name": "not-a-sensor",
		  "lat": -36.8779565276809,
		  "lon": 174.7881226266269744,
		  "tags": [
			"a",
			"b",
			"c"
		  ]
		}
	`)

	// Should return a 404
	require.Equal(t, http.StatusNotFound, rr.Code)
	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"error": "no sensor resource exists: not-a-sensor",
	}, res)
}

func TestUpdateSensorByName_Invalid(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor to work with using POST /sensors
	rr := httpRequest(t, router, "POST", "/sensors", `
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

	// Update the sensor, with invalid json values
	rr = httpRequest(t, router, "PUT", "/sensors/abc123", `
		{
		  "not": "valid",
		  "sensor": "data"
		}
	`)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"error": "invalid request body: json: unknown field \"not\"",
	}, res)
}

func TestUpdateSensorByName_NewName(t *testing.T) {
	// TODO: define this behavior...?
	require.True(t, false, "TODO")
}

func createSensor(t *testing.T, router *SensorRouter, sensorJSON string) *httptest.ResponseRecorder {
	// TODO: could remove this method, extra unneeded abstraction
	return httpRequest(t, router, "POST", "/sensors", sensorJSON)
}

func httpRequest(t *testing.T, router *SensorRouter, method string, url string, body string) *httptest.ResponseRecorder {
	handler := router.Handler()
	rr := httptest.NewRecorder()

	// Prepare body as an io.Reader (if supplied)
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	// Prepare request
	req, err := http.NewRequest(method, url, bodyReader)
	require.NoError(t, err)

	// Set application/json header
	req.Header.Set("Content-Type", "application/json")

	// Send request
	handler.ServeHTTP(rr, req)

	return rr
}

func unmarshalResponseJSON(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	var res map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &res)
	require.NoError(t, err)

	return res
}
