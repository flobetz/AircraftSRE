package main

import (
	"net/http"
	"log"
	"encoding/json"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/julienschmidt/httprouter"
	"time"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/goji/httpauth"
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
var (
	err error
	aircraft1 = Aircraft{Name: "DHC-8-400", SeatCount: 80}
	aircraft2 = Aircraft{Name: "Boeing B737", SeatCount: 186}
	aircraft3 = Aircraft{Name: "Airbus A340", SeatCount: 300}
	dbAlive = promauto.NewGauge(prometheus.GaugeOpts{
		Name:      "db_is_alive",
		Help:      "0 = DB is dead, 1 = DB is up and running",
	})
	amountFlightsGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name:		"app_flights_counter",
		Help: 		"counter of DELETE endpoint /v1/flights/<flightnumber>",
	})
	createCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name:		"app_create_flight_counter",
		Help: 		"counter of POST endpoint /v1/flights",
	})
	getAllCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name:		"app_get_all_counter",
		Help: 		"counter of GET endpoint /v1/flights",
	})
	getOneCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name:		"app_get_one_counter",
		Help: 		"counter of GET endpoint /v1/flights/<flightnumber>",
	})
	deleteCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name:		"app_delete_counter",
		Help: 		"counter of DELETE endpoint /v1/flights/<flightnumber>",
	})
)

// Entrypoint
func main() {
	var (
		// define basic auth users and passwords
		user = "flightoperator"
		password = "topsecret!"
		db *sql.DB
	)

	// initialize DB
	connStr := "user=postgres dbname=postgres password=postgres host=database port=5432 sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		println("could not create connection to postgres DB")
		return
	}
	// create connection with Ping()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
		println("could not open connection to postgres DB")
		return
	}
	println("Connection to DB successful")
	db.Query("CREATE TABLE flights (flightnumber varchar(255), startloc varchar(255), endloc varchar(255), aircraft varchar(255), departure TIMESTAMP)")
	defer db.Close()

	// start serving metrics
	serveMetrics()

	// record prometheus metrics
	customMetrics(db)

	// initiate http router
	router := httprouter.New()

	// define http routes and map to functions
	router.POST("/v1/flights", BasicAuth(createFlight(db), user, password))
	router.GET("/v1/flights", BasicAuth(getAllFlight(db), user, password))
	router.GET("/v1/flights/:flightnumber", BasicAuth(getSpecificFlight(db), user, password))
	router.DELETE("/v1/flights/:flightnumber", BasicAuth(deleteFlight(db), user, password))

	// start serving API
	log.Fatal(http.ListenAndServe(":80", router))
}

func serveMetrics() {
	go func() {
		var (
			promuser = "prometheus"
			prompw = "MetriXRule!"
		)
		// defining metrics endpoint
		println("start serving metrics on port 2112")
		http.Handle("/metrics", httpauth.SimpleBasicAuth(promuser, prompw)(promhttp.Handler()))
		log.Fatal(http.ListenAndServe(":2112", nil))
	} ()
}

func customMetrics(db *sql.DB) {
	go func() {
		var amountFlights float64
		// check for db connection every 4 seconds
		for {
			err := db.QueryRow("SELECT COUNT(*) FROM  flights;").Scan(&amountFlights)
			amountFlightsGauge.Set(amountFlights)
			if err != nil {
				println("DB is not working")
				dbAlive.Set(0)
			} else {
				dbAlive.Set(1)
				println("DB is working")
			}
			time.Sleep(4 * time.Second)
		}
	}()
}

// Basic Auth for flights endpoint
func BasicAuth(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

// flights endpoint functions
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

		// increase create Counter metric
		createCounter.Inc()
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

		// increase metric counter
		getAllCounter.Inc()
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
			exists string
		)
		// get flightnumber from http request
		FlightNumber = ps.ByName("flightnumber")

		// first check db if entry exists
		err = db.QueryRow("SELECT 1 from flights WHERE flightnumber=$1;", FlightNumber).Scan(&exists)

		if exists == "1" {
			// flightnumber exists
			// get specific flight information from DB
			err = db.QueryRow("SELECT * FROM flights WHERE flightnumber=$1", FlightNumber).Scan(&FlightNumber, &StartLoc, &EndLoc, &Aircraft, &Departure)
			checkErr(err, w)

			// encode to json and send http response
			f = Flight{FlightNumber:FlightNumber, StartLoc:StartLoc, EndLoc:EndLoc, Aircraft:Aircraft, Departure:Departure}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(f)
		} else {
			// flightnumber does not exist
			w.WriteHeader(http.StatusBadRequest)
		}

		//increase metric counter
		getOneCounter.Inc()
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

		// increase metric counter
		deleteCounter.Inc()
	}
}

func checkErr(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}