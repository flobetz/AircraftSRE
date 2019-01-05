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
	"strconv"
	"github.com/pkg/errors"
)

// define used structs and interfaces
type Aircraft struct {
	Name string
	SeatCount int
}

type Flight struct {
	FlightNumber string
	Start string
	End string
	Departure time.Time
	Aircraft string
}

type NoFlightNumber struct {
	Start string
	End string
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
	db.Query("CREATE TABLE flights (flightnumber INTEGER , startloc varchar(255), endloc varchar(255), aircraft varchar(255), departure TIMESTAMP);")
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
				dbAlive.Set(0)
			} else {
				dbAlive.Set(1)
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
		// define an empty flight and a flight without flightnumber
		var (
			noNumber NoFlightNumber
			maxFlightnumber int
			newFlightnumber int
			newNumberAsString string
		)

		// validate json body
		noNumber, err = validate(request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// create new flight number
		err = db.QueryRow("SELECT COUNT(*) FROM  flights;").Scan(&maxFlightnumber)
		newFlightnumber = maxFlightnumber + 1
		newNumberAsString = strconv.Itoa(newFlightnumber)

		// paste values of json input into DB
		db.Query("INSERT INTO flights VALUES($1,$2,$3,$4,$5);", newFlightnumber, noNumber.Start, noNumber.End, noNumber.Aircraft, noNumber.Departure)
		w.WriteHeader(http.StatusOK)
		location := "Location: /v1/flights/" + newNumberAsString + "\n"
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
			Start string
			End string
			Aircraft string
			Departure time.Time
		)
		// get all flights from Database
		rows, err := db.Query("SELECT * FROM flights")
		checkErr(err, w)
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&FlightNumber, &Start, &End, &Aircraft, &Departure)
			checkErr(err, w)
			f = Flight{FlightNumber:FlightNumber, Start:Start, End:End, Aircraft:Aircraft, Departure:Departure}
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
			Start string
			End string
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
			err = db.QueryRow("SELECT * FROM flights WHERE flightnumber=$1", FlightNumber).Scan(&FlightNumber, &Start, &End, &Aircraft, &Departure)
			checkErr(err, w)

			// encode to json and send http response
			f = Flight{FlightNumber:FlightNumber, Start:Start, End:End, Aircraft:Aircraft, Departure:Departure}
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

func validate(request *http.Request) (NoFlightNumber, error) {
	noNumber := NoFlightNumber{}

	err = json.NewDecoder(request.Body).Decode(&noNumber)

	// check incoming body for content
	if request.Body == nil {
		err = errors.New("Please send a JSON body")
		return noNumber, err
	}

	if err != nil {
		err = errors.New("JSON body does not fit")
		return noNumber, err
	}

	// check if start empty
	if noNumber.Start == "" {
		err = errors.New("Field start has to be given")
		return noNumber, err
	}

	// check if end empty
	if noNumber.End == "" {
		err = errors.New("Field end has to be given")
		return noNumber, err
	}

	// check if aircraft empty
	if noNumber.Aircraft == "" {
		err = errors.New("Field aircraft has to be given")
		return noNumber, err
	}

	// check if departure empty
	if noNumber.Departure.String() == "" {
		err = errors.New("Field departure has to be given")
		return noNumber, err
	}

	// check if departure time is in future
	if noNumber.Departure.Before(time.Now()) {
		err = errors.New("Departure time must be in future")
		return noNumber, err
	}

	// check aircraftname of request
	if ! (noNumber.Aircraft == aircraft1.Name || noNumber.Aircraft == aircraft2.Name || noNumber.Aircraft == aircraft3.Name) {
		err = errors.New("Aircraft type is not correct, please choose 'DHC-8-400', 'Boeing B737' or 'Airbus A340'")
		return noNumber, err
	}

	err = nil
	return noNumber, err
}