package main

import (
	"net/http"
	"log"
	"encoding/json"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/julienschmidt/httprouter"
	"time"
)

// define used structs and interfaces
type Aircraft struct {
	Name string
	SeatCount int
}

type Flight struct {
	FlightNumber string
	StartLoc string
	EndLoc string
	Departure time.Time
	Aircraft string
}

// define global variables and constants
var db *sql.DB
var err error
var aircraft1 = Aircraft{Name: "DHC-8-400", SeatCount: 80}
var aircraft2 = Aircraft{Name: "Boeing B737", SeatCount: 186}
var aircraft3 = Aircraft{Name: "Airbus A340", SeatCount: 300}


// Entrypoint
func main() {
	// create ai
	// init db
	// establish DB connection
	connStr := "user=postgres dbname=postgres password=postgres host=database port=5432 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		println("could not create connection to postgres DB")
		return
	}

	if err = db.Ping(); err != nil {
		println(err)
		println("could not open connection to postgres DB")
		return
	}

	println("Connection to DB successful")
	// this should be done on startup of the postgres docker container
	db.Query("CREATE TABLE flights (flightnumber varchar(255), startloc varchar(255), endloc varchar(255), aircraft varchar(255), departure TIMESTAMP)")

	defer db.Close()

	// initiate http router
	router := httprouter.New()

	// define http routes and map to functions
	router.POST("/v1/flights", createFlight(db))
	router.GET("/v1/flights", getAllFlight(db))
	router.GET("/v1/flights/:flightnumber", getSpecificFlight(db))
	router.DELETE("/v1/flights/:flightnumber", deleteFlight(db))

	// start serving API
	log.Fatal(http.ListenAndServe(":80", router))
}

func createFlight(db *sql.DB) func(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		// define an empty flight
		var f Flight
		// check incoming body for content
		if request.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		// json decode request body apply to the flight object, save fields in flight f and output in err
		err := json.NewDecoder(request.Body).Decode(&f)
		// check err
		checkErr(err, w)

		// check aircraftname of request
		if ! (f.Aircraft == aircraft1.Name || f.Aircraft == aircraft2.Name || f.Aircraft == aircraft3.Name) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// paste values of json input into DB
		db.Query("INSERT INTO flights VALUES($1,$2,$3,$4,$5);", f.FlightNumber, f.StartLoc, f.EndLoc, f.Aircraft, f.Departure)
		w.WriteHeader(http.StatusOK)
		location := "Location: /v1/flights/" + f.FlightNumber + "\n"
		w.Write([]byte(location))
	}
}


func getAllFlight(db *sql.DB) func(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		var (
			flights []Flight
			f Flight
			FlightNumber string
			StartLoc string
			EndLoc string
			Aircraft string
			Departure time.Time
		)
		// get all flights from Database
		rows, err := db.Query("SELECT * FROM flights")
		checkErr(err, w)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&FlightNumber, &StartLoc, &EndLoc, &Aircraft, &Departure)
			checkErr(err, w)
			f = Flight{FlightNumber:FlightNumber, StartLoc:StartLoc, EndLoc:EndLoc, Aircraft:Aircraft, Departure:Departure}
			flights = append(flights, f)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(flights)
	}
}

func getSpecificFlight(db *sql.DB) func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		var (
			f Flight
			FlightNumber string
			StartLoc string
			EndLoc string
			Aircraft string
			Departure time.Time
		)
		// get flightnumber from http request
		FlightNumber = ps.ByName("flightnumber")

		// get specific flight information from DB
		err = db.QueryRow("SELECT * FROM flights WHERE flightnumber=$1", FlightNumber).Scan(&FlightNumber, &StartLoc, &EndLoc, &Aircraft, &Departure)
		checkErr(err, w)

		// encode to json and send http response
		f = Flight{FlightNumber:FlightNumber, StartLoc:StartLoc, EndLoc:EndLoc, Aircraft:Aircraft, Departure:Departure}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(f)
	}
}

func deleteFlight(db *sql.DB) func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	return func(w http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		var (
			FlightNumber string
		)
		FlightNumber = ps.ByName("flightnumber")

		// delete flight from db
		db.Query("DELETE FROM flights WHERE flightnumber=$1", FlightNumber)
		w.WriteHeader(http.StatusOK)
	}
}

func checkErr(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
}