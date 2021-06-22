package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
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

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// loop until main route responds
	weather.WaitForServerRoute(s.Addr + "/weather")

	// create channel to wait for CTRL + C
	// and prevent main from exiting
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	s.Shutdown()
}
