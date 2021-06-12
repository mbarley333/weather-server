package server

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
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

type WeathersHandlers struct {
	Store map[string]Weather
}

func (h *WeathersHandlers) Get(w http.ResponseWriter, r *http.Request) {
	sliceWeather := make([]Weather, len(h.Store))

	i := 0
	for _, weather := range h.Store {
		sliceWeather[i] = weather
		i++
	}

	//marshal json for Write
	jsonBytes, err := json.Marshal(sliceWeather)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func WeatherReport(w http.ResponseWriter, req *http.Request) {
	f, err := os.Open("testdata/weather_test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	w.WriteHeader(http.StatusOK)
	io.Copy(w, f)

}

func StartServer() {

	var wait time.Duration
	var h = WeathersHandlers{}

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	r.HandleFunc("/weatherreport", WeatherReport)
	r.HandleFunc("/weathers", h.Get)

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
