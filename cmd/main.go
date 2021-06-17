package main

import (
	"flag"
	"log"
	"weather"
)

func main() {

	// Config struct
	config := new(weather.Config)

	// server startup configs
	flag.IntVar(&config.Port, "port", 9010, "HTTP Server port")
	flag.StringVar(&config.LogLevel, "log", "verbose", "log level: verbose or quiet")

	flag.Parse()

	s := config.NewServer()
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
