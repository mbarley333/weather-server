package main

import (
	"log"
	"weather"
)

func main() {

	s := weather.NewServer()
	log.Fatal(s.ListenAndServe())

}
