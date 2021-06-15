package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// struct for json
type Weather struct {
	Main        string  `json:"main"`
	Description string  `json:"description"`
	Temp        float64 `json:"temp"`
	City        string  `json:"city"`
}

func (s *server) handleWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	// get actual weather
	conditions, err := s.GetWeather(city)
	if err != nil {
		msg := fmt.Sprintf("unable to locate city %q", city)
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

type server struct {
	httpServer *http.Server
	Addr       string
	GetWeather func(string) (Weather, error)
}

func (s *server) ListenAndServe() error {
	// Add a 'quiet' flag to disable logging?
	// log.Println("Starting up on ", s.Addr)
	s.httpServer = &http.Server{
		Addr:              s.Addr,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/weather", s.handleWeather)
	s.httpServer.Handler = mux

	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Println("server start:", err)
		return err
	}
	return nil
}

func (s *server) Shutdown() {
	// //create channel for os.Signal comands
	// c := make(chan os.Signal, 1)

	// //accept CTRL+C for granceful shutdown, unblock condition
	// signal.Notify(c, os.Interrupt)

	// //everything below the <-c is the activity once channel unblocks
	// <-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s.httpServer.Shutdown(ctx)
}

func GetWeatherFromOpenWeatherMap(city string) (Weather, error) {
	// call OWM
	return Weather{City: "Not implemented yet"}, nil
}

func NewServer(port int) server {
	return server{
		Addr:       fmt.Sprintf("127.0.0.1:%d", port),
		GetWeather: GetWeatherFromOpenWeatherMap,
	}
}
