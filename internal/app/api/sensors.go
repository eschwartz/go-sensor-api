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
		return nil, 400, fmt.Errorf("invalid request body: %w", err)
	}

	// Store the new sensor
	createdSensor, err := router.store.Create(sensor)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to store sensor: %w", err)
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
		return nil, http.StatusInternalServerError, err
	}

	// Handle no matching sensor
	if sensor == nil {
		return nil, http.StatusNotFound, fmt.Errorf("no sensor exists with name \"%s\"", name)
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
