package main

import (
	"fmt"
	"net/http"
	server "weather-server"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", server.Hello)
	fmt.Println("Starting up on :9000")
	http.ListenAndServe(":9000", r)
}
