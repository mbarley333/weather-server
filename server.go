package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"weather/api"
)

// struct to config new server
type Config struct {
	Port     int
	LogLevel string
}

// struct for query string parameters to pass
// to OWM api
type UrlParameters struct {
	City  string
	Units string
}

// struct for json
type Weather struct {
	Main        string  `json:"main"`
	Description string  `json:"description"`
	Temp        float64 `json:"temp"`
	City        string  `json:"city"`
}

// since server has state (start, shutdown, etc) use a struct
// to hold the object
// server struct hold necesssary info to start up a weather HTTP server
// GetWeather accepts a func with the following signature: string input and return Weather, error
type server struct {
	httpServer *http.Server
	Addr       string
	GetWeather func(UrlParameters) (Weather, error)
	logger     *log.Logger
}

// NewServer creates HTTP and is useful to fulfill parallel testing
// by passing in a different port # per test
func (c Config) NewServer() server {

	newLogger := log.New(os.Stdout, "", log.LstdFlags)
	if c.LogLevel == "quiet" {
		newLogger.SetOutput(ioutil.Discard)
	}

	// new server with address:port
	// and assigns the GetWeatherFromOpenWeatherMap method to the GetWeather field
	return server{
		Addr:       fmt.Sprintf("127.0.0.1:%d", c.Port),
		GetWeather: GetWeatherFromOpenWeatherMap,
		logger:     newLogger,
	}
}

// ListenAndServe starts up the HTTP server with
// specific paramters to override the default settings
// also needed to allow for multiple instances of HTTP server
// to run in parallel without port conflicts.  The Addr field
// needs to be set with the unique port
func (s *server) ListenAndServe() error {

	s.httpServer = &http.Server{
		Addr:              s.Addr,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
		ErrorLog:          s.logger,
	}

	s.logger.Println("Starting up on ", s.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/weather", s.handleWeather)
	s.httpServer.Handler = mux

	if err := s.httpServer.ListenAndServe(); err != nil {
		s.logger.Println("server start:", err)
		return err
	}

	return nil
}

// Shutdown is a method on the server struct
// and performs HTTP server shutdown
func (s *server) Shutdown() {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.logger.Println("shutting down..")
	s.httpServer.Shutdown(ctx)
	os.Exit(0)

}

func (s *server) handleWeather(w http.ResponseWriter, r *http.Request) {

	// url paramater logic -- possibly move to func
	params := UrlParameters{
		City:  r.URL.Query().Get("city"),
		Units: r.URL.Query().Get("units"),
	}

	if params.Units == "" {
		params.Units = "imperial"
	}

	// get actual weather
	conditions, err := s.GetWeather(params)
	if err != nil {
		msg := fmt.Sprintf("unable to locate city %q", params.City)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	// marshal json for Write
	jsonBytes, err := json.Marshal(conditions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// GetWeatherFromOpenWeatherMap will access to OWM api to pull in weather data
func GetWeatherFromOpenWeatherMap(params UrlParameters) (Weather, error) {
	// call OWM

	apiKey, err := api.GetWeatherAPIKey("WEATHERAPI")
	if err != nil {
		log.Fatal("Unable to get API key")
	}

	client, err := api.NewClient(apiKey, params.Units)
	if err != nil {
		log.Fatal("Something went wrong")
	}

	response, err := client.Get(params.City)
	if err != nil {
		log.Println(err)
	}

	return Weather(response), nil
}

// waitForServerRoute checks if the main route is reachable
func WaitForServerRoute(url string) {

	for {
		_, err := net.Dial("tcp", url)
		if err == nil {
			log.Println("tcp not listening")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}

}
