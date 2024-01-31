package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eschwartz/pingthings-sensor-api/internal/app/store"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type SensorRouter struct {
	store store.SensorStore
}

func NewSensorRouter() *SensorRouter {
	return &SensorRouter{
		// TODO: Replace with a persistent data store
		store: store.NewMemorySensorStore(),
		// NEXT: and api endpoints for CRUD using store, test
	}
}

func (router *SensorRouter) Handler() http.Handler {
	r := mux.NewRouter()

	// GET /health - Health Check
	r.HandleFunc("/health", WithJSONHandler(router.HealthCheckHandler)).
		Methods("GET")

	// POST /sensors - Create Sensor
	r.HandleFunc("/sensors", WithJSONHandler(router.CreateSensorHandler)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	// GET /sensors/{name} - Get Sensor by Name
	r.HandleFunc("/sensors/{name}", WithJSONHandler(router.GetSensorByNameHandler)).
		Methods("GET")

	// PUT /sensors/{name} - Update Sensor by Name
	r.HandleFunc("/sensors/{name}", WithJSONHandler(router.UpdateSensorByNameHandler)).
		Methods("PUT")

	return r
}

func (router *SensorRouter) HealthCheckHandler(r *http.Request) (interface{}, int, error) {
	return map[string]bool{
		"ok": true,
	}, http.StatusOK, nil
}

func (router *SensorRouter) CreateSensorHandler(r *http.Request) (interface{}, int, error) {
	// Parse JSON request body
	sensor, err := decodeSensorJSON(r.Body)
	if err != nil {
		// Errors are mostly likely caused by malformed request bodies
		// Here's a nice write-up on decoder error handling, if we want something
		// more precise: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		return nil, 400, err
	}

	// Store the new sensor
	createdSensor, err := router.store.Create(sensor)
	if err != nil {
		// Unknown error from store, log and respond as 500
		log.Printf("failed to create sensor: %s", err)
		return nil, 500, fmt.Errorf("failed to store sensor: %w", errors.New("internal server error"))
	}

	return SensorDetailsResponse{*createdSensor}, http.StatusCreated, nil
}

func (router *SensorRouter) GetSensorByNameHandler(r *http.Request) (interface{}, int, error) {
	// Get sensor {name} from URL
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		// Missing {name} means we probably misconfigured the route
		log.Println("GET /sensors/{name} request is missing the \"name\" var.")
		return nil, http.StatusInternalServerError, errors.New("interval server error")
	}

	// Retrieve sensor from data store
	sensor, err := router.store.GetByName(name)
	if err != nil {
		// Unknown error from store, log and respond as 500
		log.Printf("failed to retrieve sensor by name \"%s\": %s", name, err)
		return nil, http.StatusInternalServerError, errors.New("interval server error")
	}

	// Handle no matching sensor
	if sensor == nil {
		return nil, http.StatusNotFound, &store.MissingResourceError{name, "sensor"}
	}

	return SensorDetailsResponse{*sensor}, http.StatusOK, nil
}

func (router *SensorRouter) UpdateSensorByNameHandler(r *http.Request) (interface{}, int, error) {
	// Get sensor {name} from URL
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		// Missing {name} means we probably misconfigured the route
		log.Println("PUT /sensors/{name} request is missing the \"name\" var.")
		return nil, http.StatusInternalServerError, errors.New("interval server error")
	}

	// Decode sensor JSON body
	sensor, err := decodeSensorJSON(r.Body)
	if err != nil {
		// Invalid request body, respond with 400
		return nil, http.StatusBadRequest, err
	}

	// Update the sensor in the data store
	sensor, err = router.store.UpdateByName(name, sensor)
	if err != nil {
		// If there's not matching resource, return a 404
		var missingErr *store.MissingResourceError
		if errors.As(err, &missingErr) {
			return nil, http.StatusNotFound, err
		}

		// Any other errors are treated as 500s
		log.Printf("failed to update store in PUT /sensors/%s: %s", name, err)
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}

	return SensorDetailsResponse{*sensor}, http.StatusOK, nil
}

func decodeSensorJSON(r io.Reader) (*store.Sensor, error) {
	// Parse JSON request body
	decoder := json.NewDecoder(r)
	// don't allow extra fields in request body
	decoder.DisallowUnknownFields()

	var sensor store.Sensor
	err := decoder.Decode(&sensor)
	if err != nil {
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	return &sensor, nil
}

type SensorDetailsResponse struct {
	Data store.Sensor `json:"data"`
}
