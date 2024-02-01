package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/eschwartz/pingthings-sensor-api/internal/app/store"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	router := NewSensorRouter()

	// Send GET /health request
	rr := httpRequest(t, router, "GET", "/health", "")
	require.Equal(t, http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, true, res["ok"])
}

func TestCreateSensor(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor using POST /sensors
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

	// Check json response
	res := unmarshalResponseJSON(t, rr)
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
	router := NewSensorRouter()

	// Create a sensor with an invalid payload
	rr := httpRequest(t, router, "POST", "/sensors", `
		{
		  "not": "valid",
		  "sensor": "data"
		}
	`)
	// Should respond with a 400
	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateSensor_StoreFailure(t *testing.T) {
	// Use FailingSensorStore, to test
	// the behavior of the API when the storage backend fails
	router := &SensorRouter{
		store: &FailingSensorStore{},
	}

	// Attempt to create a sensor, with a failing store
	rr := httpRequest(t, router, "POST", "/sensors", `
		{
		  "name": "abc123",
		  "lat": 44.916241209323736,
		  "lon": -93.21112681214602,
		  "tags": []
		}
	`)
	// Should respond with a 500
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	// Should respond with an error message
	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"error": "failed to store sensor: internal server error",
	}, res)
}

func TestGetSensorByName(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor using POST /sensors
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

	// Get the sensor, using GET /sensors/:name
	rr = httpRequest(t, router, "GET", "/sensors/abc123", "")
	require.Equal(t, http.StatusOK, rr.Code)

	// Should return the sensor that we created earlier
	res := unmarshalResponseJSON(t, rr)
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
	rr := httpRequest(t, router, "GET", "/sensors/not-a-sensor", "")

	// Should respond with a 404
	require.Equal(t, http.StatusNotFound, rr.Code)

	// Should include an error message
	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"error": "no sensor resource exists: not-a-sensor",
	}, res)
}

func TestGetSensor_StoreFailure(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor using POST /sensors
	rr := httpRequest(t, router, "POST", "/sensors", `
		{
		  "name": "abc123",
		  "lat": 44.916241209323736,
		  "lon": -93.21112681214602,
		  "tags": []
		}
	`)
	// Should response w/201
	require.Equal(t, http.StatusCreated, rr.Code)

	// Use a failing store backend, to test how the GET endpoint handles it
	router.store = &FailingSensorStore{}

	// Get the sensor, using GET /sensors/:name
	rr = httpRequest(t, router, "GET", "/sensors/abc123", "")
	// Should return a 500
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	// Should return the sensor that we created earlier
	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"error": "failed to retrieve sensor: interval server error",
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

func TestUpdateSensor_StoreFailure(t *testing.T) {
	router := NewSensorRouter()

	// Create sensor using POST /sensors
	rr := httpRequest(t, router, "POST", "/sensors", `
		{
		  "name": "abc123",
		  "lat": 44.916241209323736,
		  "lon": -93.21112681214602,
		  "tags": []
		}
	`)
	// Should response w/201
	require.Equal(t, http.StatusCreated, rr.Code)

	// Use a failing store backend, to test how the GET endpoint handles it
	router.store = &FailingSensorStore{}

	// Update the sensor, using PUT /sensors/:name
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
	// Should return a 500
	require.Equal(t, http.StatusInternalServerError, rr.Code)

	// Should return the sensor that we created earlier
	res := unmarshalResponseJSON(t, rr)
	require.Equal(t, map[string]interface{}{
		"error": "failed to update sensor: internal server error",
	}, res)
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

// FailingSensorStore is an implementation of SensorStore,
// where every method returns an error.
// This may be used to integration test API behavior in case of failing store backends
type FailingSensorStore struct {
}

func (f FailingSensorStore) Create(sensor *store.Sensor) (*store.Sensor, error) {
	return nil, errors.New("FailingSensorStore.Create() failing for tests, on purpose")
}

func (f FailingSensorStore) GetByName(name string) (*store.Sensor, error) {
	return nil, errors.New("FailingSensorStore.GetByName() failing for tests, on purpose")
}

func (f FailingSensorStore) UpdateByName(name string, sensor *store.Sensor) (*store.Sensor, error) {
	return nil, errors.New("FailingSensorStore.UpdateByName() failing for tests, on purpose")
}
