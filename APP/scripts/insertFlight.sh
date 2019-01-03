#!/usr/bin/env bash
curl -v --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"2","StartLoc":"2","EndLoc":"2","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://40.80.144.36:80/v1/flights