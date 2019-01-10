#!/usr/bin/env bash
set -x
URL=http://flightoperator.eastus.cloudapp.azure.com

# insert flights
printf "insert flight 1:\n"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"DHC-8-400","Departure":"2020-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "insert flight 2"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"Boeing B737","Departure":"2020-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "insert flight 3"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"Airbus A340","Departure":"2020-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

# get all flights
printf "get all flights:"
curl -i -u flightoperator:topsecret! --request GET \
  $URL/v1/flights
printf "\n\n"
sleep 7

# get only one flight
printf "get only one flight:"
curl -i -u flightoperator:topsecret! --request GET \
  $URL/v1/flights/1
printf "\n\n"
sleep 7

# delete flight one
printf "delete one flight"
curl -i -u flightoperator:topsecret! --request DELETE \
  $URL/v1/flights/3
printf "\n\n"
sleep 7

# get all flights again
printf "get all flights agin:"
curl -i -u flightoperator:topsecret! --request GET \
  $URL/v1/flights
printf "\n\n"
sleep 7

# test basic auth
printf "wrongpw:"
curl -i -u flightoperator:wrongpw --request GET \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "wronguser:"
curl -i -u wronguser:topsecret! --request GET \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "no auth:"
curl -i --request GET \
  $URL/v1/flights
printf "\n\n"
sleep 7

# testing http responses
printf "testing http responses:"
printf ""
printf "successful DB entry:"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"3","End":"3","Aircraft":"DHC-8-400","Departure":"2020-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "unsuccessful DB entry:"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"flightNumber":"1"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "successful get All flights"
curl -i -u flightoperator:topsecret! --request GET \
    $URL/v1/flights
printf "\n\n"
sleep 7

printf "unsuccessful get All flights"
curl -i -u flightoperator:topsecret! --request GET \
  $URL/v1/flight
printf "\n\n"
sleep 7

printf "successful get specific flight"
curl -i -u flightoperator:topsecret! --request GET \
  $URL/v1/flights/1
printf "\n\n"
sleep 7

printf "unsuccessful get specific flight"
curl -i -u flightoperator:topsecret! --request GET \
  $URL/v1/flights/doesnotexist
printf "\n\n"
sleep 7

printf "successful delete flight"
curl -i -u flightoperator:topsecret! --request DELETE \
  $URL/v1/flights/3
printf "\n\n"
sleep 7

printf "unsuccessful delete flight"
curl -i -u flightoperator:topsecret! --request DELETE \
  $URL/v1/flights/doesnotexist
printf "\n\n"
sleep 7

printf "try to insert flight with wrong aircraft:"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"doesnotexist","Departure":"2020-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "try to insert flight with departure time in past:"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"Start":"1","End":"1","Aircraft":"Airbus A340","Departure":"2015-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7

printf "try to insert flight with wrong json format"
curl -i -u flightoperator:topsecret! --header "Content-Type: application/json" \
  --request POST \
  --data '{"fieldDoesNotExist":"1","End":"1","Aircraft":"Airbus A340","Departure":"2015-10-10T10:00:00Z"}' \
  $URL/v1/flights
printf "\n\n"
sleep 7