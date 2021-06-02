package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
