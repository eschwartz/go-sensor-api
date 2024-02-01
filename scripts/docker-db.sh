#!/usr/bin/env bash

docker run \
    -p 5432:5432 \
    -e POSTGRES_USER=admin \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -e POSTGRES_DB=sensors \
    -v ./scripts/db-init.sql:/docker-entrypoint-initdb.d/db-init.sql \
    -d \
    postgis/postgis
