package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eschwartz/go-sensor-api/internal/app/store"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

// Regexp for parsing radius query parameters
// eg "50km"
var radiusParamRegexp = regexp.MustCompile("([0-9]+)(km|mi)")

// Regexp for parsing location query parameter
// eg 45.12,-90.34
var latLonRegexp = regexp.MustCompile("^(-?[0-9]+\\.?[0-9]*),(-?[0-9]+\\.?[0-9]*)$")

type SensorRouter struct {
	store store.SensorStore
}

func NewSensorRouter() (*SensorRouter, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, errors.New("must set DATABASE_URL")
	}

	postgisStore, err := store.NewPostgisStore(dbUrl)
	if err != nil {
		return nil, err
	}

	return &SensorRouter{
		store: postgisStore,
	}, nil
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

	// GET /sensors/closest?location=&radius=
	r.HandleFunc("/sensors/closest", WithJSONHandler(router.FindClosestSensor)).
		Queries("location", "{location}", "radius", "{radius}")

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
		return nil, http.StatusInternalServerError, errors.New("failed to retrieve sensor: interval server error")
	}

	// Handle no matching sensor
	if sensor == nil {
		return nil, http.StatusNotFound, &store.MissingResourceError{name, "sensor"}
	}

	return SensorDetailsResponse{*sensor}, http.StatusOK, nil
}

func (router *SensorRouter) FindClosestSensor(r *http.Request) (interface{}, int, error) {
	// Load query params
	vars := mux.Vars(r)
	locationParam, ok := vars["location"]
	if !ok {
		return nil, http.StatusBadRequest, errors.New("missing required \"location\" param")
	}
	radiusParam, ok := vars["radius"]
	if !ok {
		return nil, http.StatusBadRequest, errors.New("missing required \"radius\" param")
	}

	// Parse radius, eg "50km"
	radiusMatch := radiusParamRegexp.FindStringSubmatch(radiusParam)
	if radiusMatch == nil {
		return nil, http.StatusBadRequest,
			errors.New("invalid value for \"radius\": must be formatted like \"50km\" or \"100mi\"")
	}
	// If the regex matches, we should always have 2 groups. If not, we didn't something wrong here
	if len(radiusMatch) != 3 {
		log.Printf("GET /sensors/closest: Unexpected number of regexp match groups for radius: \"%s\"", radiusParam)
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}
	// Convert radius to meters
	radiusValue, err := strconv.Atoi(radiusMatch[1])
	if err != nil {
		return nil, http.StatusBadRequest,
			errors.New("invalid value for \"radius\": must be formatted like \"50km\" or \"100mi\"")
	}
	radiusUnits := radiusMatch[2]
	var radiusMeters int
	if radiusUnits == "km" {
		radiusMeters = radiusValue * 1000
	} else if radiusUnits == "mi" {
		radiusMeters = int(float64(radiusValue) * 1609.34)
	} else {
		return nil, http.StatusBadRequest,
			errors.New("invalid unit for \"radius\": must be \"km\" or \"mi\"")
	}

	// Parse location, eg 45.12,-90.34
	locationMatch := latLonRegexp.FindStringSubmatch(locationParam)
	if locationMatch == nil {
		// Attempt to geocode the location, assuming it's a place name / address

		return nil, http.StatusBadRequest,
			errors.New("invalid value for \"location\": must be formatted like \"45.12,-90.34")
	}
	// If the regex matches, we should always have 2 groups. If not, we didn't something wrong here
	if len(locationMatch) != 3 {
		log.Printf("GET /sensors/closest: Unexpected number of regexp match groups for location: \"%s\"", locationParam)
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}
	lat, err := strconv.ParseFloat(locationMatch[1], 64)
	if err != nil {
		return nil, http.StatusBadRequest,
			errors.New("invalid value for \"location\": must be formatted like \"45.12,-90.34")
	}
	lon, err := strconv.ParseFloat(locationMatch[2], 64)
	if err != nil {
		return nil, http.StatusBadRequest,
			errors.New("invalid value for \"location\": must be formatted like \"45.12,-90.34")
	}

	// Lookup closest sensors
	sensors, err := router.store.FindClosest(lat, lon, radiusMeters)
	if err != nil {
		log.Printf("GET /sensors/closest failed to FindClosest(): %s", err)
		return nil, http.StatusInternalServerError, errors.New("internal server error")
	}

	return SensorListResponse{sensors}, http.StatusOK, nil
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
		return nil, http.StatusInternalServerError, errors.New("failed to update sensor: internal server error")
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

type SensorListResponse struct {
	Data []*store.Sensor `json:"data"`
}
