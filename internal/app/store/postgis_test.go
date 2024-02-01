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

func TestPostgisStore_CreateAndGetByName(t *testing.T) {
	// Skip tests unless the test DB env var is set
	dbUrl := os.Getenv("TEST_DATABASE_URL")
	if dbUrl == "" {
		t.Skip("Skipping database tests")
	}

	// Reset DB before and after each test
	resetDb(t, dbUrl)
	defer resetDb(t, dbUrl)

	store, err := NewPostgisStore(dbUrl)
	require.NoError(t, err)
	defer store.Close()

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
	// Skip tests unless the test DB env var is set
	dbUrl := os.Getenv("TEST_DATABASE_URL")
	if dbUrl == "" {
		t.Skip("Skipping database tests")
	}

	// Reset DB before and after each test
	resetDb(t, dbUrl)
	defer resetDb(t, dbUrl)

	store, err := NewPostgisStore(dbUrl)
	require.NoError(t, err)
	defer store.Close()

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

func TestNewPostgisStore_GetByNameMissing(t *testing.T) {
	require.True(t, false)
}

func TestPostgisStore_UpdateByName(t *testing.T) {
	// Skip tests unless the test DB env var is set
	dbUrl := os.Getenv("TEST_DATABASE_URL")
	if dbUrl == "" {
		t.Skip("Skipping database tests")
	}

	// Reset DB before and after each test
	resetDb(t, dbUrl)
	defer resetDb(t, dbUrl)

	store, err := NewPostgisStore(dbUrl)
	require.NoError(t, err)
	defer store.Close()

	// Create a sensor
	_, err = store.Create(&Sensor{
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

func TestNewPostgisStore_UpdateMissing(t *testing.T) {
	require.True(t, false)
}
