#!/usr/bin/env bash
# insert flights
echo "insert flight 1"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"DHC-8-400","Departure":"2020-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "insert flight 2"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"Boeing B737","Departure":"2020-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "insert flight 3"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"Airbus A340","Departure":"2020-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

# get all flights
echo "get all flights:"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

# get only one flight
echo "get only one flight:"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/1
echo ""
echo ""
sleep 3

# delete flight one
echo "delete one flight"
echo ""
curl -i -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/3
echo ""
echo ""
sleep 3

# get all flights again
echo "get all flights agin:"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

# test basic auth
echo "wrongpw:"
echo ""
curl -i -u flightoperator:wrongpw --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "wronguser:"
echo ""
curl -i -u wronguser:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "no auth:"
echo ""
curl -i --request GET \
  http://localhost:80/v1/flights
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
  --data '{"Start":"3","End":"3","Aircraft":"DHC-8-400","Departure":"2020-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
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
sleep 3

echo "successful get All flights"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "unsuccessful get All flights"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flight
echo ""
echo ""
sleep 3

echo "successful get specific flight"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/1
echo ""
echo ""
sleep 3

echo "unsuccessful get specific flight"
echo ""
curl -i -u flightoperator:topsecret! --request GET \
  http://localhost:80/v1/flights/doesnotexist
echo ""
echo ""
sleep 3

echo "successful delete flight"
echo ""
curl -i -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/3
echo ""
echo ""
sleep 3

echo "unsuccessful delete flight"
echo ""
curl -i -u flightoperator:topsecret! --request DELETE \
  http://localhost:80/v1/flights/doesnotexist
echo ""
echo ""
sleep 3

echo "try to insert flight with wrong aircraft:"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"doesnotexist","Departure":"2020-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "try to insert flight with departure time in past:"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"Airbus A340","Departure":"2015-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3

echo "try to insert flight with wrong json format"
echo ""
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"fieldDoesNotExist":"1","End":"1","Aircraft":"Airbus A340","Departure":"2015-10-10T10:00:00Z"}' \
  http://localhost:80/v1/flights
echo ""
echo ""
sleep 3