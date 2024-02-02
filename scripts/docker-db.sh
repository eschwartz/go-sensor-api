#!/usr/bin/env bash
# Start a postgres docker container, with the postgis extension installed
# The SQL code in ./db-init.sql will be run on startup
#
# The database may be accessed using the following connection URL
# postgres://admin:mysecretpassword@localhost:5432/sensors?sslmode=disable
docker run \
    -p 5432:5432 \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -e POSTGRES_DB=sensors \
    -v ./scripts/db-init.sql:/docker-entrypoint-initdb.d/db-init.sql \
    -d \
    postgis/postgis
