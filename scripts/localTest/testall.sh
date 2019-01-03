#!/usr/bin/env bash
# insert flights
echo "insert flight 1"
echo ""
curl -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"1","StartLoc":"1","EndLoc":"1","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "insert flight 2"
echo ""
curl -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"2","StartLoc":"2","EndLoc":"2","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "insert flight 3"
echo ""
curl -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"3","StartLoc":"3","EndLoc":"3","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

# get all flights
echo "get all flights:"
echo ""
curl -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

# get only one flight
echo "get only one flight:"
echo ""
curl -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/1
echo ""
echo ""
echo ""
sleep 3

# delete flight one
echo "delete one flight"
echo ""
curl -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/1
echo ""
echo ""
echo ""
sleep 3

# get all flights again
echo "get all flights agin:"
echo ""
curl -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

# test basic auth
echo "wrongpw:"
echo ""
curl -u flightoperator:wrongpw --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "wronguser:"
echo ""
curl -u wronguser:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "no auth:"
echo ""
curl --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

# testing http responses
echo "testing http responses:"
echo ""
echo "successful DB entry:"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"1","StartLoc":"1","EndLoc":"1","Aircraft":"DHC-8-400","Departure":"2018-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "unsuccessful DB entry:"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"1"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "successful get All flights"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
echo ""
sleep 3

echo "unsuccessful get All flights"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flight
echo ""
echo ""
echo ""
sleep 3

echo "successful get specific flight"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/1
echo ""
echo ""
echo ""
sleep 3

echo "unsuccessful get specific flight"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/doesnotexist
echo ""
echo ""
echo ""
sleep 3

echo "successful delete flight"
echo ""
curl -i -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/1
echo ""
echo ""
echo ""
sleep 3

echo "unsuccessful delete flight"
echo ""
curl -i -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/doesnotexist

