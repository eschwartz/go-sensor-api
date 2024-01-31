package api

import (
	"encoding/json"
	"fmt"
	"github.com/eschwartz/pingthings-sensor-api/internal/app/store"
	"github.com/gorilla/mux"
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

	r.HandleFunc("/health", WithJSONHandler(router.HealthCheckHandler))
	r.HandleFunc("/sensors", WithJSONHandler(router.CreateSensorHandler)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	return r
}

func (router *SensorRouter) HealthCheckHandler(r *http.Request) (interface{}, int, error) {
	return map[string]bool{
		"ok": true,
	}, 200, nil
}

func (router *SensorRouter) CreateSensorHandler(r *http.Request) (interface{}, int, error) {
	// Parse JSON request body
	// TODO: abstract for reuse
	decoder := json.NewDecoder(r.Body)
	// don't allow extra fields in request body
	decoder.DisallowUnknownFields()

	var sensor store.Sensor
	err := decoder.Decode(&sensor)
	if err != nil {
		// Errors are mostly likely caused by malformed request bodies
		// Here's a nice write-up on decoder error handling, if we want something
		// more precise: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
		return nil, 400, fmt.Errorf("invalid request body: %w", err)
	}

	// Store the new sensor
	createdSensor, err := router.store.Create(&sensor)
	if err != nil {
		return nil, 500, fmt.Errorf("failed to store sensor: %w", err)
	}

	return SensorResponse{*createdSensor}, 201, nil
}

type SensorResponse struct {
	Data store.Sensor `json:"data"`
}
