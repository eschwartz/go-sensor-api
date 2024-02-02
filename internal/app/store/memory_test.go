package store

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	store := NewMemorySensorStore()

	sensor := &Sensor{
		Name: "abc123",
		Lat:  10,
		Lon:  20,
		Tags: []string{"a", "b"},
	}

	// Create the sensor resource
	createdSensor, err := store.Create(sensor)
	require.NoError(t, err)
	assert.Same(t, sensor, createdSensor)
}

func TestGetByName(t *testing.T) {
	store := NewMemorySensorStore()

	sensor := &Sensor{
		Name: "abc123",
		Lat:  10,
		Lon:  20,
		Tags: []string{"a", "b"},
	}

	// Create the sensor resource
	_, err := store.Create(sensor)
	require.NoError(t, err)

	// Retrieve the sensor by name
	retrievedSensor, err := store.GetByName("abc123")
	require.NoError(t, err)
	assert.Same(t, sensor, retrievedSensor)
}

func TestUpdate(t *testing.T) {
	store := NewMemorySensorStore()

	sensor := &Sensor{
		Name: "abc123",
		Lat:  10,
		Lon:  20,
		Tags: []string{"a", "b"},
	}

	// Create the sensor resource
	_, err := store.Create(sensor)
	require.NoError(t, err)

	// Update the sensor
	newSensor := &Sensor{
		Name: "abc123",
		Lat:  5,
		Lon:  7,
		Tags: []string{"x", "y"},
	}
	updatedSensor, err := store.UpdateByName("abc123", newSensor)
	require.NoError(t, err)
	require.Same(t, newSensor, updatedSensor)
}

func TestUpdateMissing(t *testing.T) {
	store := NewMemorySensorStore()

	// Attempt to update a sensor that does not exit
	newSensor := &Sensor{
		Name: "abc123",
		Lat:  5,
		Lon:  7,
		Tags: []string{"x", "y"},
	}
	updatedSensor, err := store.UpdateByName("abc123", newSensor)
	require.Nil(t, updatedSensor)
	require.NotNil(t, err)
	require.IsType(t, err, &MissingResourceError{})
	require.Equal(t, "no sensor resource exists: abc123", err.Error())
}
