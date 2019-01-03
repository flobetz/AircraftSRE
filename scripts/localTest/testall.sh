#!/usr/bin/env bash
# insert flights
echo "insert flight 1"
curl -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"1","StartLoc":"1","EndLoc":"1","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights

echo "insert flight 2"
curl -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"2","StartLoc":"2","EndLoc":"2","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights

echo "insert flight 3"
curl -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"3","StartLoc":"3","EndLoc":"3","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights

# get all flights
echo "get all flights:"
curl -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights

# get only one flight
echo "get only one flight:"
curl -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/1

# delete flight one
echo "delete one flight"
curl -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/1

# get all flights again
echo "get all flights agin:"
curl -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights

# test basic auth
echo "wrongpw:"
curl -u flightoperator:wrongpw --request GET \
  http://localhost:80/v1/flights

echo "wronguser:"
curl -u wronguser:topsecret! --request GET \
  http://localhost:80/v1/flights

echo "no auth:"
curl --request GET \
  http://localhost:80/v1/flights