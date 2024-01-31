package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type SensorRouter struct {
}

func NewSensorRouter() *SensorRouter {
	return &SensorRouter{}
}

func (router *SensorRouter) Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/health", router.HealthCheckHandler)

	return r
}

func (router *SensorRouter) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	if err != nil {
		log.Printf("failed to encode json response: %s", err)
		w.WriteHeader(500)
	}
}
