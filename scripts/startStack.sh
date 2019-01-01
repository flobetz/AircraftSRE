#!/usr/bin/env bash

# start postgres container
docker run --name some-postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres

# compile go app
cd /Users/flobetz/Documents/projects/go/AircraftSRE
go build

# run go app
./AircraftSRE