package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	UnitsMetric   = "metric"
	UnitsImperial = "imperial"
	UnitsStandard = "standard"
)

var units = map[string]bool{
	UnitsImperial: true,
	UnitsMetric:   true,
	UnitsStandard: true,
}

// WeatherReponse is used to accept JSON structured
// data.  Output is not very human friendly and is thus
// a stage type for the Weather struct
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

// Client is used to assemble the necessary parts for a HTTP request
// and includes the HTTPClient lib to set timeout
type Client struct {
	Base       string
	Units      string
	APIKey     string
	HTTPClient *http.Client
}

// Get takes a location and resturns a Weather struct and error.
func (c Client) Get(location string) (Weather, error) {

	//assemble request based on Client struct data
	url := c.FormatURL(location)

	//use Client since we override default HTTPClient settings -- timeout
	resp, err := c.HTTPClient.Get(url)

	if err != nil {
		return Weather{}, fmt.Errorf("error contacting api: %v", err)
	}
	//close when done to prevent resource leaks
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Weather{}, fmt.Errorf("unexpected status code.  %v", resp.StatusCode)
	}

	// pass io.Reader to ParseReponse
	return ParseResponse(resp.Body)

}

// NewClient takes apiKey and tempunits populates the Client struct
// with a base url, temperature units, apikey AND sets up the
// HTTPClient with timeout settings since default timeout is too long.
// Returns Client struct and error
func NewClient(apiKey string, tempunits string) (Client, error) {

	//var c Client

	result := validUnit(tempunits)
	if !result {
		return Client{}, fmt.Errorf("invalid unit of measurement: %s", tempunits)
	}

	c := Client{
		Units:      tempunits,
		Base:       "https://api.openweathermap.org",
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	return c, nil

}

// GetWeatherAPIKey takes env as name of environmental variable
// which holds the API key
func GetWeatherAPIKey(env string) (string, error) {

	apikey := os.Getenv(env)

	if apikey == "" {
		return "", fmt.Errorf("%s value not set", env)
	}
	return apikey, nil
}

// FormatURL is a method on the Client struct that
// assembles the URL used in the request
func (c Client) FormatURL(location string) string {

	return fmt.Sprintf("%s/data/2.5/weather?q=%s&units=%s&appid=%s", c.Base, location, c.Units, c.ApiKey)

}

// ParseReponse takes an io.Reader and decodes into WeatherReponse struct
// The WeatherResonse struct in then used to setup the Weather struct for
// human reading
func ParseResponse(r io.Reader) (Weather, error) {

	var result WeatherResponse

	//decodes io.Reader into variable address (e.g. &result)
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return Weather{}, err
	}

	w := Weather{
		Main:        result.Weather[0].Main,
		Description: result.Weather[0].Description,
		Temp:        result.Main.Temp,
		City:        result.City,
	}

	return w, nil
}

func validUnit(unit string) bool {
	return units[unit]
}
