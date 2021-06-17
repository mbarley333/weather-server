package main

import (
	"flag"
	"log"
	"weather"
)

func main() {

	//Config struct
	config := new(weather.Config)

	//assign cli flags to Config struct fields
	flag.IntVar(&config.Port, "port", 9010, "HTTP Server port")
	flag.StringVar(&config.LogLevel, "log", "quiet", "log level: verbose or quiet")

	flag.Parse()

	s := config.NewServer()
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
