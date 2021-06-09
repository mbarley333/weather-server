package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func Hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
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
	r := mux.NewRouter()
	r.HandleFunc("/hello", Hello)
	r.HandleFunc("/weatherreport", WeatherReport)

	srv := &http.Server{
		Addr:              "127.0.0.1:9000",
		Handler:           r,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}
	fmt.Println("Starting up on :9000")

	log.Fatal(srv.ListenAndServeTLS(os.Getenv("CERTPATH"), os.Getenv("KEYPATH")))

}
