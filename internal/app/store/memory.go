package store

// MemorySensorStore is an in-memory store of Sensor models
// Future iterations should transition to a persistent data store
// (though the in-memory may continue to be useful for testing)
type MemorySensorStore struct {
	byName map[string]*Sensor
}

func NewMemorySensorStore() *MemorySensorStore {
	return &MemorySensorStore{
		byName: make(map[string]*Sensor),
	}
}

func (s *MemorySensorStore) Create(sensor *Sensor) (*Sensor, error) {
	// TODO: validate sensor input
	// TODO: validate unique name
	s.byName[sensor.Name] = sensor

	return sensor, nil
}

func (s *MemorySensorStore) GetByName(name string) (*Sensor, error) {
	sensor, ok := s.byName[name]
	if !ok {
		return nil, nil
	}

	return sensor, nil
}

func (s *MemorySensorStore) UpdateByName(name string, sensor *Sensor) (*Sensor, error) {
	_, ok := s.byName[name]
	if !ok {
		return nil, &MissingResourceError{
			id:           name,
			resourceType: "sensor",
		}
	}

	// TODO validate sensor data
	// TODO what if name has changed? Is this allowed?
	s.byName[name] = sensor

	return sensor, nil
}
