package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type MapboxGeoService struct {
	http              *http.Client
	mapboxAccessToken string
}

func (svc *MapboxGeoService) Geocode(place string) (float64, float64, error) {
	// Call mapbox API to geocode the place name
	// into lat/lon coordinates
	placeNameEncoded := url.PathEscape(place)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.mapbox.com/geocoding/v5/mapbox.places/%s.json", placeNameEncoded),
		nil,
	)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("failed to prepare Mapbox request: %w", err)
	}

	// Add the Mapbox access token
	q := req.URL.Query()
	q.Add("access_token", svc.mapboxAccessToken)
	req.URL.RawQuery = q.Encode()

	// Execute request
	resp, err := svc.http.Do(req)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("request to Mapbox Geocode service failed: %w", err)
	}

	// TODO check status code in response

	var geocodeResp MapboxGeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&geocodeResp)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode mapbox geocode response: %w", err)
	}

	// Place not found
	if len(geocodeResp.Features) == 0 {
		return 0, 0, fmt.Errorf("no location found at %s", place)
	}
	// Expect feature to have a center
	if len(geocodeResp.Features[0].Center) != 2 {
		return 0, 0, errors.New("mapbox geocode response has invalid center")
	}

	return geocodeResp.Features[0].Center[1], geocodeResp.Features[0].Center[0], nil
}

type MapboxGeocodeResponse struct {
	Features []MapboxGeocodeResponseFeature `json:"features"`
}

type MapboxGeocodeResponseFeature struct {
	Center []float64 `json:"center"`
}
