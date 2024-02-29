# Sample Go API: Sensor Metadata

A sample REST API written in Go. Covers a hypothetical use case that tracks metadata for physical sensors in various geographical locations.

## Features

A JSON REST API for storing and querying sensor metadata.
Exposes endpoints for the following:

- Storing name, location (gps position), and a list of tags for each sensor.
- Retrieving metadata for an individual sensor by name.
- Updating a sensor’s metadata.
- Querying to find the sensor nearest to a given location (by lat/lon).
- Query to find sensor nearest to a location by place name (geocoded).


## Usage

Install dependencies

```
go mod tidy
```

Setup a local database

```sh
# See "Database Setup" below
./scripts/docker-db.sh
export DATABASE_URL=postgres://admin:mysecretpassword@localhost:5432/sensors?sslmode=disable
```

Run the HTTP server

```
go run ./cmd/sensor-api
```

This will start the API server on `localhost:8000`

## Database Setup

The sensors API requires a Postgres database with the postgis extension installed. To create the necessary database tables, see [`scripts/db-init.sql`](./scripts/db-init.sql) 

For local development, you can quickly create a test database in a docker container:

```sh
./scripts/docker-db.sh
```

This will start a `postgis/postgis` docker container initialized with the necessary database tables. You may connect to that database using:

```
postgres://admin:mysecretpassword@localhost:5432/sensors?sslmode=disable
```

## Tests

To run tests:

```
go test ./...
```

Some tests require a test database, and will be skipped if none is specified. To specify the test database, use the `TEST_DATABASE_URL` environment variable:

```sh
# Create a test database in a docker container
./scripts/docker-db.sh

# Configure the database url
export TEST_DATABASE_URL=postgres://admin:mysecretpassword@localhost:5432/sensors?sslmode=disable

# Run tests
go test ./...
```

> Test database tables **will be truncated** with every test run. **Do not use a live database!**


### Environment configuration

The following environment variables are supported:

| Name         | Description                                |
|--------------|--------------------------------------------|
| PORT         | HTTP port to listen on. Defaults to `8000` |
| DATABASE_URL | URL to connect to the postgres database    |


## API Reference

### GET /sensors/:name

Retrieve metadata for a single sensor, by name.

#### Example

```
GET /sensors/abc123
```

```json
HTTP 200
{
    "data": {
      "id": 1234,
      "name": "abc123",
      "lat": 44.916241209323736,
      "lon": -93.21112681214602,
      "tags": [
        "x",
        "y",
        "z"
      ]
    }
}
```

### GET /sensors/closest

Retrieve metadata for sensors closest to a given location.

#### Example

```
GET /sensors/closest/?location=44.9,-93.211&radius=100km
```

```json
HTTP 200
{
    "data": [
      {
        "id": 1234,
        "name": "abc123",
        "lat": 44.916241209323736,
        "lon": -93.21112681214602,
        "tags": [
          "x",
          "y",
          "z"
        ]
      }
    ]
}
```


#### Query Parameters

| Parameter | Required | Default | Description                                                                                                             | Example         |
|-----------|----------|---------|-------------------------------------------------------------------------------------------------------------------------|-----------------|
| location  | x        | -       | Latitude / longitute coordinate, from which to center the search                                                        | `44.9,-93.211`  |
| radius    |          | `20km`  | Results will be included within this radius from the `location`. Supported units are `mi` (miles) and `km` (kilometers) | `50mi`, `100km` |           

### POST /sensors

Add a sensor to the system

#### Example

```json
POST /sensors
{
  "name": "abc123",
  "lat": 44.916241209323736,
  "lon": -93.21112681214602,
  "tags": [
    "x",
    "y",
    "z"
  ]
}
```


```json
HTTP 201
{
    "data": {
      "id": 1234,
      "name": "abc123",
      "lat": 44.916241209323736,
      "lon": -93.21112681214602,
      "tags": [
        "x",
        "y",
        "z"
      ]
    }
}
```

### PUT /sensors/:name

Update a sensor's metadata, by sensor name

#### Example

```json
PUT /sensors/abc123
{
  "name": "abc123",
  "lat": -36.8779565276809,
  "lon": 174.7881226266269744,
  "tags": [
    "a",
    "b",
    "c"
  ]
}
```

```json
HTTP 200
{
    "data": {
      "id": 1234,
      "name": "abc123",
      "lat": -36.8779565276809,
      "lon": 174.7881226266269744,
      "tags": [
        "a",
        "b",
        "c"
      ]
    }
}
```
