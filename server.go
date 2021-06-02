package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WeatherResponse struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	City string `json:"name"`
}

// Weather is the human friendly struct that is populated
// by the WeatherReponse struct
type Weather struct {
	Main        string
	Description string
	Temp        float64
	City        string
}

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

// ParseReponse takes an io.Reader and decodes into WeatherReponse struct
// The WeatherResonse struct in then used to setup the Weather struct for
// human reading
func ParseResponse(r io.Reader) (Weather, error) {

	var result WeatherResponse

	// decodes io.Reader into variable address (e.g. &result)
	// binds the json to the struct
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return Weather{}, err
	}

	// human friendly struct
	var w Weather

	w.Main = result.Weather[0].Main
	w.Description = result.Weather[0].Description
	w.Temp = result.Main.Temp
	w.City = result.City
	return w, nil
}
