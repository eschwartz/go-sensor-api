package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/cridenour/go-postgis"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

type PostgisStore struct {
	db *sql.DB
}

func NewPostgisStore(dbUrl string) (*PostgisStore, error) {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &PostgisStore{
		db: db,
	}, nil
}

func (store *PostgisStore) Create(sensor *Sensor) (*Sensor, error) {
	// Begin the DB transaction
	tx, err := store.db.BeginTx(context.Background(), nil)
	defer tx.Rollback()

	// Insert the sensor record
	createSql := `
		INSERT INTO sensors (name, location) 
		-- see https://postgis.net/docs/ST_MakePoint.html
		--VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326))
		VALUES ($1, GeomFromEWKB($2))
		RETURNING id;
	`
	var id int
	err = tx.QueryRow(createSql, sensor.Name, newGisPoint(sensor.Lat, sensor.Lon)).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	// Update the sensor ID
	sensor.ID = id

	// Insert tags
	err = store.createSensorTags(sensor.ID, sensor.Tags, tx)
	if err != nil {
		return nil, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return sensor, nil
}

func (store *PostgisStore) GetByName(name string) (*Sensor, error) {
	sql := `
		SELECT 
			sensors.id, 
			sensors.location,
			-- Join in tags, as a nested array
			array_remove(array_agg(tags.value), NULL) as tags
		FROM sensors
		LEFT JOIN tags on sensors.id = tags.sensor_id
		WHERE sensors.name = $1
		GROUP BY sensors.id
	`
	var id int
	location := newGisPoint(0, 0)
	var tags pq.StringArray
	err := store.db.QueryRow(sql, name).
		Scan(&id, &location, &tags)
	if err != nil {
		// We want to return nil if there are no matches
		// sql lib does not have typed errors, so we need to match on a string here
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &Sensor{
		ID:   id,
		Name: name,
		Lon:  location.X,
		Lat:  location.Y,
		Tags: tags,
	}, nil
}

func (store *PostgisStore) UpdateByName(name string, sensor *Sensor) (*Sensor, error) {
	// Begin the DB transaction
	tx, err := store.db.BeginTx(context.Background(), nil)
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(`
		UPDATE sensors
		SET name = $2, location = GeomFromEWKB($3)
		WHERE name = $1
		RETURNING id
	`, name, sensor.Name, newGisPoint(sensor.Lat, sensor.Lon)).Scan(&id)
	if err != nil {
		// Handle no match errors
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("failed to update sensor: no sensor exists with name \"%s\"", name)
		}
		return nil, err
	}

	sensor.ID = id

	// Replace all the tags
	// TODO: There's probably a way to do this that avoids unnecessary deletion
	// Delete all the tags....
	_, err = store.db.Exec(`
		DELETE FROM tags
		WHERE sensor_id = $1
	`, sensor.ID)
	if err != nil {
		return nil, err
	}
	// ...then recreate them all
	if err := store.createSensorTags(sensor.ID, sensor.Tags, tx); err != nil {
		return nil, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return sensor, nil
}

func (store *PostgisStore) FindClosest(lat float64, lon float64, radiusMeters int) ([]*Sensor, error) {
	// Query DB for closest sensors
	rows, err := store.db.Query(`
		SELECT 
			sensors.id, 
			sensors.name,
			sensors.location,
			-- Join in tags, as a nested array
			array_remove(array_agg(tags.value), NULL) as tags,
			ST_Distance(sensors.location::geography, GeomFromEWKB($1)::geography) as distance
		FROM sensors
		LEFT JOIN tags on sensors.id = tags.sensor_id
		-- find within radius
		--WHERE ST_DWithin(sensors.location, ST_MakePoint($1, $2)::geography, $3)
		WHERE ST_DWithin(sensors.location::geography, GeomFromEWKB($1)::geography, $2)
		GROUP BY sensors.id
		-- sort by distance
		ORDER BY ST_Distance(sensors.location::geography, GeomFromEWKB($1)::geography);
	`, newGisPoint(lat, lon), radiusMeters)
	if err != nil {
		return []*Sensor{}, err
	}
	defer rows.Close()

	// Iterate through results, to create slice of Sensors
	var sensors []*Sensor
	for rows.Next() {
		// hydrate values from DB row
		var id int
		var name string
		var tags pq.StringArray
		var distance float64
		location := newGisPoint(0, 0)
		if err := rows.Scan(&id, &name, &location, &tags, &distance); err != nil {
			return sensors, nil
		}

		// Create a sensor for db row data
		sensors = append(sensors, &Sensor{
			ID:   id,
			Name: name,
			Lon:  location.X,
			Lat:  location.Y,
			Tags: tags,
		})
	}

	return sensors, rows.Err()
}

func (store *PostgisStore) Close() error {
	return store.db.Close()
}

func (store *PostgisStore) createSensorTags(sensorId int, tags []string, tx *sql.Tx) error {
	if len(tags) == 0 {
		return nil
	}

	tagSql := `
			INSERT INTO tags (sensor_id, value)
			VALUES 
		`
	// Need to dynamically create VALUES, to insert a row for every specified tag
	var tagValuesSqls []string
	tagSqlArgs := []any{sensorId}
	for i, tag := range tags {
		tagValuesSqls = append(tagValuesSqls, fmt.Sprintf("($1, $%d)", i+2))
		tagSqlArgs = append(tagSqlArgs, tag)
	}
	tagSql += strings.Join(tagValuesSqls, ", ")

	_, err := tx.Exec(tagSql, tagSqlArgs...)
	return err
}

func newGisPoint(lat, lon float64) *postgis.PointS {
	return &postgis.PointS{SRID: 4326, X: lon, Y: lat}
}
