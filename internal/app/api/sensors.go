package api

import (
	"github.com/eschwartz/pingthings-sensor-api/internal/app/store"
	"github.com/gorilla/mux"
	"net/http"
)

type SensorRouter struct {
	store store.SensorStore
}

func NewSensorRouter() *SensorRouter {
	return &SensorRouter{}
}

func (router *SensorRouter) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/health", WithJSONHandler(router.HealthCheckHandler))

	return r
}

func (router *SensorRouter) HealthCheckHandler(r *http.Request) (interface{}, int, error) {
	return map[string]bool{
		"ok": true,
	}, 200, nil
}
