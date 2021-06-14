package server

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

//struct for json
type Weather struct {
	Id          string  `json:"id"`
	Main        string  `json:"main"`
	Description string  `json:"description"`
	Temp        float64 `json:"temp"`
	City        string  `json:"city"`
}

type WeatherHandlers struct {
	sync.Mutex
	Store map[string]Weather
}

func (h *WeatherHandlers) weather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetByCity(w, r)
		return
	case "POST":
		h.Post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return

	}

}

func (h *WeatherHandlers) GetByCity(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filters := strings.ToLower(query.Get("city"))

	var sliceWeather []Weather

	h.Lock()

	for _, weather := range h.Store {
		if strings.ToLower(weather.City) == filters {
			sliceWeather = append(sliceWeather, weather)
		}
	}
	h.Unlock()

	if len(sliceWeather) == 0 {
		w.Write([]byte("unable to locate city"))
		return
	}

	//marshal json for Write
	jsonBytes, err := json.Marshal(sliceWeather)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *WeatherHandlers) Get(w http.ResponseWriter, r *http.Request) {
	sliceWeather := make([]Weather, len(h.Store))
	h.Lock()
	i := 0
	for _, weather := range h.Store {
		sliceWeather[i] = weather
		i++
	}
	h.Unlock()

	//marshal json for Write
	jsonBytes, err := json.Marshal(sliceWeather)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *WeatherHandlers) Post(w http.ResponseWriter, r *http.Request) {

	var result Weather
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer r.Body.Close()

	h.Lock()
	h.Store[result.Id] = result
	defer h.Unlock()
}

func StartServer() {

	var wait time.Duration
	var h = WeatherHandlers{
		Store: map[string]Weather{
			"id1": {
				Id:          "id1",
				Main:        "Cloudy",
				Description: "Partly cloudy",
				Temp:        74.6,
				City:        "Kaneohe",
			},
			"id2": {
				Id:          "id2",
				Main:        "Rain",
				Description: "Passing showers",
				Temp:        64.6,
				City:        "Seattle",
			},
		},
	}

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	//r.HandleFunc("/weatherreport", WeatherReport)
	r.HandleFunc("/weather", h.weather)

	srv := &http.Server{
		Addr:              "127.0.0.1:9000",
		Handler:           r,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}
	fmt.Println("Starting up on :9000")

	// Run our server in a goroutine so that it doesn't block.
	// https://github.com/gorilla/mux
	go func() {
		if err := srv.ListenAndServeTLS(os.Getenv("CERTPATH"), os.Getenv("KEYPATH")); err != nil {
			log.Println(err)
		}
	}()

	//create channel the understand the os.Signal comands
	c := make(chan os.Signal, 1)

	//accept CTRL+c for granceful shutdown, unblock condition
	signal.Notify(c, os.Interrupt)

	//everything below the <-c is the activity once channel unblocks
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)

}
