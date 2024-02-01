package store

import "fmt"

type MissingResourceError struct {
	ID           string
	ResourceType string
}

func (e *MissingResourceError) Error() string {
	return fmt.Sprintf("no %s resource exists: %s", e.ResourceType, e.ID)
}

type Sensor struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Lat  float64  `json:"lat"`
	Lon  float64  `json:"lon"`
	Tags []string `json:"tags"`
}

type SensorStore interface {
	Create(sensor *Sensor) (*Sensor, error)
	GetByName(name string) (*Sensor, error)
	UpdateByName(name string, sensor *Sensor) (*Sensor, error)
}
