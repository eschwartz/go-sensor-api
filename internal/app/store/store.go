package store

import "fmt"

type MissingResourceError struct {
	id           string
	resourceType string
}

func (e *MissingResourceError) Error() string {
	return fmt.Sprintf("No %s resource exists: %s", e.resourceType, e.id)
}

type Sensor struct {
	Name string
	Lat  float64
	Lon  float64
	Tags []string
}

type SensorStore interface {
	Create(sensor *Sensor) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	UpdateByName(name string, sensor *Sensor) (*Sensor, error)
}
