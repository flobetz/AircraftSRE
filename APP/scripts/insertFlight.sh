#!/usr/bin/env bash
curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"1","StartLoc":"1","EndLoc":"1","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights