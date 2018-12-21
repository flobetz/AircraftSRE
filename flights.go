package main

import (
	"net/http"
	"log"
	"encoding/json"
	"fmt"
	//"strconv"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/davecgh/go-spew/spew"
)

// define global variables and constants
var flightDB []Flight
var requestedFlight int
var db sql.DB

// define used structs and interfaces
type Aircraft struct {
	Name string
	SeatCount int
}

type Flight struct {
	FlightNumber int
	StartLoc string
	EndLoc string
	Aircraft Aircraft
}

// Entrypoint
func main() {
	// establish DB connection
	connStr := "user=postgres dbname=postgres password=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		println("could not create connection to postgres DB")
	}

	if err = db.Ping(); err != nil {
		err = errors.Wrapf(err,
			"Couldn't ping postgre database (%s)",
			spew.Sdump(connStr))
		println(err)
		println("could not open connection to postgres DB")
	}

	// this should be done on startup of the postgres docker container
	db.Query("CREATE TABLE flights (flightnumber varchar(255), startloc varchar(255), endloc varchar(255), aircraft varchar(255))")
	db.Query("CREATE TABLE planes (name varchar(255), seatcount int)")


    // print flightDB array first
	printFlightArray(flightDB)

	// initiate http router
	router := mux.NewRouter()

	// define http routes and map to functions
	router.HandleFunc("/v1/flights", createFlight).Methods("POST")
	router.HandleFunc("/v1/flights", getAllFlight).Methods("GET")
	router.HandleFunc("/v1/flights/{id}", getSpecificFlight).Methods("GET")
	router.HandleFunc("/v1/flights/{id}", deleteFlight).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":80", router))
}

func createFlight(w http.ResponseWriter, request *http.Request) {
	db.Query("INSERT INTO flights VALUES ('flug1', 'ehingen', 'stuttgart', 'flugzeug1')")
	//db.Prepare("INSERT INTO flights(flightnumber, startloc, endloc, aircraft) VALUES (flug1, ehingen, stuttgart, flugzeug1)")
	w.Write([]byte("flight saved!\n"))
}

func deleteFlight(w http.ResponseWriter, request *http.Request) {
	db.Query("DELETE FROM flights WHERE flightnumber='flug1'")
	w.Write([]byte("flight deleted\n"))
}

func getAllFlight(w http.ResponseWriter, request *http.Request) {
	// get all flights from Database
	db.Query("SELECT * FROM flights")
	fmt.Fprint(w, "here are all saved flights:\n")
	json, _ := json.Marshal(flightDB)
	w.Write(json)
}

func getSpecificFlight(w http.ResponseWriter, request *http.Request) {
	w.Write([]byte("Get specific flight\n"))
	vars := mux.Vars(request)
	id, err := vars["id"]
	if err != true {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", id)


	// get specific flight information from DB
	//w.Write([]byte("Quering DB"))
	//db.Query("SELECT * FROM flights WHERE flightnumber=?", id)
}

func printFlightArray(flightDB []Flight) {
	fmt.Printf("len=%d cap=%d %v\n", len(flightDB), cap(flightDB), flightDB)
}