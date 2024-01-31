package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSONHandlerFunc is an alternate signature for http handler function
// The function will return some response data, a status code, and (optionally) an error
type JSONHandlerFunc func(r *http.Request) (interface{}, int, error)

// WithJSONHandler converts a JSONHandlerFunc to a standard http.HandlerFun
// This allows for standardized serialization of response data and error
func WithJSONHandler(f JSONHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the underlying handler function
		data, status, httpErr := f(r)

		// Serve handler errors as JSON
		if httpErr != nil {
			data = map[string]string{
				"error": httpErr.Error(),
			}
		}

		// Write JSON response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(data)

		// Handle JSON encoding failure
		if err != nil {
			log.Printf("failed to encode json response: %s", err)
			w.WriteHeader(500)
		}
	}
}
