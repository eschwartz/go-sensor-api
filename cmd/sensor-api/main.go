package main

import (
	"fmt"
	"github.com/eschwartz/pingthings-sensor-api/internal/app/api"
	"log"
	"net/http"
	"os"
)

func main() {
	router := api.NewSensorRouter()

	// Read port from env var, or use default val
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// HTTP Listen
	log.Printf("Listening on http://localhost:%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router.Handler())
	log.Fatal(err)
}
