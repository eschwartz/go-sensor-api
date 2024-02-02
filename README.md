# Sensor Metadata API

Technical challenge for PingThings Go Developer/SRE application, submitted by Edan Schwartz, Feb 2024.

## Task

Please build a JSON REST API for storing and querying sensor metadata.
At a minimum, this API should expose endpoints for the following:

- Storing name, location (gps position), and a list of tags for each sensor.
- Retrieving metadata for an individual sensor by name.
- Updating a sensor’s metadata.
- Querying to find the sensor nearest to a given location.

It is up to you how you structure your application, but please write it in Go and include anything you would
in a professional project (i.e.: README, tests, input validation, etc).

## Usage

Install dependencies

```
go mod tidy
```

Run the HTTP server

```
go run ./cmd/sensor-api
```

### Environment configuration

The following environment variables are supported:

| Name | Description                                |
|------|--------------------------------------------|
| PORT | HTTP port to listen on. Defaults to `8000` |


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
      "id": "18595f26-d6c6-4fe8-999e-77f87d9edad0",
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
        "id": "18595f26-d6c6-4fe8-999e-77f87d9edad0",
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
      "id": "18595f26-d6c6-4fe8-999e-77f87d9edad0",
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
      "id": "18595f26-d6c6-4fe8-999e-77f87d9edad0",
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
