package store

import (
	"database/sql"
	"fmt"
	"github.com/cridenour/go-postgis"
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

	// Insert the sensor record
	createSql := `
		INSERT INTO sensors (name, location) 
		-- see https://postgis.net/docs/ST_MakePoint.html
		--VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326))
		VALUES ($1, GeomFromEWKB($2))
		RETURNING id;
	`
	var id int
	err := store.db.QueryRow(createSql, sensor.Name, postgis.PointS{SRID: 4326, X: sensor.Lon, Y: sensor.Lat}).
		Scan(&id)
	if err != nil {
		return nil, err
	}

	// Update the sensor ID
	sensor.ID = id

	// Insert tags
	// TODO wrap this all in a transaction, so we don't get dangling sensors without tags
	if len(sensor.Tags) > 0 {
		tagSql := `
			INSERT INTO tags (sensor_id, value)
			VALUES 
		`
		// Need to dynamically create VALUES, to insert a row for every specified tag
		var tagValuesSqls []string
		tagSqlArgs := []any{sensor.ID}
		for i, tag := range sensor.Tags {
			tagValuesSqls = append(tagValuesSqls, fmt.Sprintf("($1, $%d)", i+2))
			tagSqlArgs = append(tagSqlArgs, tag)
		}
		tagSql += strings.Join(tagValuesSqls, ", ")

		_, err = store.db.Exec(tagSql, tagSqlArgs...)
		if err != nil {
			return nil, err
		}
	}

	return sensor, nil
}

func (store *PostgisStore) GetByName(name string) (*Sensor, error) {
	//TODO implement me
	panic("implement me")
}

func (store *PostgisStore) UpdateByName(name string, sensor *Sensor) (*Sensor, error) {
	//TODO implement me
	panic("implement me")
}

func (store *PostgisStore) Close() error {
	return store.db.Close()
}
