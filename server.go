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
	"os/signal"
	"syscall"
	"time"
)

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
	Units       string  `json:"units"`
}

// since server has state (start, shutdown, etc) use a struct
// to hold the object
// server struct hold necesssary info to start up a weather HTTP server
// GetWeather accepts a func with the following signature: string input and return Weather, error
type Server struct {
	httpServer *http.Server
	Addr       string
	GetWeather func(UrlParameters) (Weather, error)
	logger     *log.Logger
	Port       int
	LogLevel   string
}

// type to hold options for Server struct
type Option func(*Server)

// override default port for Server
// provides cleaner user experience
func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithLogLevel(loglevel string) Option {
	return func(s *Server) {
		s.LogLevel = loglevel
	}
}

// NewServer creates HTTP and is useful to fulfill parallel testing
// by passing in a different port # per test.  Takes in a variadic
// parameter, opts, to use a Server options
func NewServer(opts ...Option) *Server {

	// create Server instance with defaults
	s := &Server{
		Port:     9090,
		LogLevel: "verbose",
	}

	// set override options.  loop takes in
	// With funcs loaded with input params and
	// executes to update Server struct
	for _, o := range opts {
		o(s)
	}

	newLogger := log.New(os.Stdout, "", log.LstdFlags)
	if s.LogLevel == "quiet" {
		newLogger.SetOutput(ioutil.Discard)
	}

	// update struct...perhaps there is a better way.
	s.Addr = fmt.Sprintf("127.0.0.1:%d", s.Port)
	s.GetWeather = GetWeatherFromOpenWeatherMap
	s.logger = newLogger

	return s

}

// ListenAndServe starts up the HTTP server with
// specific paramters to override the default settings
// also needed to allow for multiple instances of HTTP server
// to run in parallel without port conflicts.  The Addr field
// needs to be set with the unique port
func (s *Server) ListenAndServe() error {

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
		WaitForServerRoute(s.Addr + "/weather")
		s.logger.Println("server start:", err)
		return err
	}

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c
	s.Shutdown()

	return nil
}

// Shutdown is a method on the server struct
// and performs HTTP server shutdown
func (s *Server) Shutdown() {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	s.logger.Println("shutting down..")
	s.httpServer.Shutdown(ctx)
	os.Exit(0)

}

func (s *Server) handleWeather(w http.ResponseWriter, r *http.Request) {

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

	apiKey, err := GetWeatherAPIKey("WEATHERAPI")
	if err != nil {
		log.Fatal("Unable to get API key")
	}

	client, err := NewClient(apiKey, params.Units)
	if err != nil {
		log.Fatalf("unable to create new client: %s", err)
	}

	response, err := client.Get(params.City)
	if err != nil {
		log.Printf("invalid city requested: %s", params.City)
	}

	// add units to Weather struct
	response.Units = params.Units

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
