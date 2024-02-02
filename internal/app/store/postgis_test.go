package store

import (
	"database/sql"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func resetDb(t *testing.T, dbUrl string) {
	db, err := sql.Open("postgres", dbUrl)
	require.NoError(t, err)

	_, err = db.Exec(`
		TRUNCATE sensors CASCADE;
		TRUNCATE tags;
	`)
	require.NoError(t, err)
}

func testSetup(t *testing.T) (*PostgisStore, func()) {
	// Skip tests unless the test DB env var is set
	dbUrl := os.Getenv("TEST_DATABASE_URL")
	if dbUrl == "" {
		t.Skip("Skipping database tests")
	}

	// Reset DB before and after each test
	resetDb(t, dbUrl)

	store, err := NewPostgisStore(dbUrl)
	require.NoError(t, err)

	// Return a cleanup function
	return store, func() {
		resetDb(t, dbUrl)
		_ = store.Close()
	}
}

func TestPostgisStore_CreateAndGetByName(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Create a sensor
	sensor, err := store.Create(&Sensor{
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
		Tags: []string{"a", "b", "c"},
	})
	require.NoError(t, err)

	require.Equal(t, &Sensor{
		ID:   sensor.ID,
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
		Tags: []string{"a", "b", "c"},
	}, sensor)
	require.NotEqual(t, 0, sensor.ID)

	// Retrieve the created sensor
	sensor, err = store.GetByName("sensor-abc")
	require.NoError(t, err)
	require.Equal(t, &Sensor{
		ID:   sensor.ID,
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
		Tags: []string{"a", "b", "c"},
	}, sensor)
	require.NotEqual(t, 0, sensor.ID)
}

func TestPostgisStore_CreateAndGetByNameNoTags(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Create a sensor with no tags
	sensor, err := store.Create(&Sensor{
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
	})
	require.NoError(t, err)

	require.Equal(t, &Sensor{
		ID:   sensor.ID,
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
	}, sensor)
	require.NotEqual(t, 0, sensor.ID)

	// Retrieve the created sensor
	sensor, err = store.GetByName("sensor-abc")
	require.NoError(t, err)
	require.Equal(t, &Sensor{
		ID:   sensor.ID,
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
		Tags: []string{},
	}, sensor)
	require.NotEqual(t, 0, sensor.ID)
}

func TestPostgisStore_GetByNameMissing(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Retrieve a sensor that doesn't exist
	sensor, err := store.GetByName("not-a-sensor")
	require.NoError(t, err)
	require.Nil(t, sensor)
}

func TestPostgisStore_UpdateByName(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Create a sensor
	_, err := store.Create(&Sensor{
		Name: "sensor-abc",
		Lat:  45.123456,
		Lon:  -90.98765,
		Tags: []string{"a", "b", "c"},
	})
	require.NoError(t, err)

	// Update the sensor
	sensor, err := store.UpdateByName("sensor-abc", &Sensor{
		Name: "sensor-xyz",
		Lat:  -36.8779565276809,
		Lon:  174.7881226266269744,
		Tags: []string{"x", "y", "z"},
	})
	require.NoError(t, err)

	// Retrieve the updated sensor
	sensor, err = store.GetByName("sensor-xyz")
	require.NoError(t, err)
	require.Equal(t, &Sensor{
		ID:   sensor.ID,
		Name: "sensor-xyz",
		Lat:  -36.8779565276809,
		Lon:  174.7881226266269744,
		Tags: []string{"x", "y", "z"},
	}, sensor)
	require.NotEqual(t, 0, sensor.ID)
}

func TestPostgisStore_UpdateMissing(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Update a sensor that does not exist
	_, err := store.UpdateByName("sensor-xyz", &Sensor{
		Name: "sensor-xyz",
		Lat:  -36.8779565276809,
		Lon:  174.7881226266269744,
		Tags: []string{"x", "y", "z"},
	})
	require.Error(t, err)
	require.Equal(t, "failed to update sensor: no sensor exists with name \"sensor-xyz\"", err.Error())
}

func TestNewPostgisStore_FindClosest(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Create sensors in multiple locations
	testSensors := []*Sensor{
		// St. Paul, MN
		{Name: "STP", Lat: 44.9558833427991, Lon: -93.09844267331863},
		// Minneapolis, MN (downtown)
		{Name: "MPLS", Lat: 44.97620767775624, Lon: -93.27360528040553},
		// Chicago, IL
		{Name: "CHI", Lat: 41.86950364771445, Lon: -87.68055283399988},
	}
	for _, sensor := range testSensors {
		_, err := store.Create(sensor)
		require.NoError(t, err)
	}

	// Find locations within 100km of S. Minneapolis
	sensors, err := store.FindClosest(44.91016213524799, -93.22412239250284, 100e3)
	require.NoError(t, err)
	require.Len(t, sensors, 2)
	require.Equal(t, Sensor{
		ID:   sensors[0].ID,
		Name: "MPLS",
		Lat:  44.97620767775624,
		Lon:  -93.27360528040553,
		Tags: []string{},
	}, *sensors[0])
	require.Equal(t, Sensor{
		ID:   sensors[1].ID,
		Name: "STP",
		Lat:  44.9558833427991,
		Lon:  -93.09844267331863,
		Tags: []string{},
	}, *sensors[1])
}

func TestNewPostgisStore_FindClosestNoResults(t *testing.T) {
	store, cleanup := testSetup(t)
	defer cleanup()

	// Find closest locations, when none exist
	sensors, err := store.FindClosest(44.91016213524799, -93.22412239250284, 100e3)
	require.NoError(t, err)

	// Should return an empty slice
	require.Len(t, sensors, 0)
}
