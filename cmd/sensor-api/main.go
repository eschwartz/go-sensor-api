package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	r := mux.NewRouter()

	// HTTP Routes
	r.HandleFunc("/health", HealthCheckHandler)

	// Read port from env var, or use default val
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// HTTP Listen
	log.Printf("Listening on http://localhost:%s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	log.Fatal(err)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	if err != nil {
		log.Printf("failed to encode json response: %s", err)
		w.WriteHeader(500)
	}
}
