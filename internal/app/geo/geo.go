package geo

type GeoService interface {
	// Returns lat/lon values, and an error
	Geocode(place string) (float64, float64, error)
}
