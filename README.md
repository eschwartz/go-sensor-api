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
