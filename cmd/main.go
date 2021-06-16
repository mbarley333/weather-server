package main

import (
	"flag"
	"log"
	"weather"
)

func main() {

	config := new(weather.Config)
	flag.IntVar(&config.Port, "port", 9010, "HTTP Server port")
	flag.StringVar(&config.LogLevel, "log", "verbose", "log level: verbose or quiet")
	flag.StringVar(&config.TempUnits, "units", "imperial", "temparature units of measurement (imperial, metric, standard)")

	// portFlag := flag.Int("port", 9010, "HTTP Server port")
	// logFlag := flag.String("log", "verbose", "log level: verbose or quiet")
	// tempUnitsFlag := flag.String("units", "imperial", "temparature units of measurement (imperial, metric, standard)")

	flag.Parse()
	s := config.NewServer()
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
