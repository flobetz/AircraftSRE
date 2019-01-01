#!/usr/bin/env bash

# start postgres container
docker run --name some-postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres

# start app container
docker run -d -p 80:80 flights:latest