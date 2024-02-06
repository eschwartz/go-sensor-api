package geo

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

func TestMapboxGeoService(t *testing.T) {
	mapboxToken := os.Getenv("MAPBOX_ACCESS_TOKEN")
	if mapboxToken == "" {
		t.Skip("Skipping MapboxGeoService live integration test. Missing MAPBOX_ACCESS_TOKEN")
	}

	svc := &MapboxGeoService{
		http:              &http.Client{},
		mapboxAccessToken: mapboxToken,
	}

	lat, lon, err := svc.Geocode("Minneapolis")
	require.NoError(t, err)
	require.NotEqual(t, 0.0, lat)
	require.NotEqual(t, 0.0, lon)
}
