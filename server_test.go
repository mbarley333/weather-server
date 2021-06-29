package weather_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"testing"
	"weather"

	"github.com/google/go-cmp/cmp"
)

// GetKaneoheTestWeather returns Kaneohe specific weather results
// and is used in the TestServerKaneohe func
func GetKaneoheTestWeather(params weather.UrlParameters) (weather.Weather, error) {
	return weather.Weather{
		Main:        "Cloudy",
		Description: "Partly cloudy",
		Temp:        74.6,
		City:        "Kaneohe",
		Units:       "imperial",
	}, nil
}

// GetSeattleTestWeather returns Seattle specific weather results
// and is used in the TestServerSeattle func
func GetSeattleTestWeather(params weather.UrlParameters) (weather.Weather, error) {
	return weather.Weather{
		Main:        "Rain",
		Description: "Passing showers",
		Temp:        64.6,
		City:        "Seattle",
		Units:       "imperial",
	}, nil
}

// GetNonexistentCityTestWeather returns an error as a result of an invalid city
func GetNonexistentCityTestWeather(params weather.UrlParameters) (weather.Weather, error) {
	return weather.Weather{}, errors.New("no weather for you")
}

// TestServerKaneohe will create an HTTP server and a GET route
// for Kaneohe
func TestServerKaneohe(t *testing.T) {
	t.Parallel()
	s := weather.NewServer(
		weather.WithPort(9000),
		weather.WithLogLevel("quiet"),
	)

	// override the GetWeatherFromOpenWeatherMap setting that prod would use
	s.GetWeather = GetKaneoheTestWeather

	//goroutine to start up HTTP server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// GET request against new HTTP server
	resp, err := http.Get("http://127.0.0.1:9000/weather?city=kaneohe")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	want := `{"main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe","units":"imperial"}`
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, string(got)) {
		t.Errorf("want %q, got %q", want, string(got))
	}
}

func TestServerNonexistentCity(t *testing.T) {
	t.Parallel()

	s := weather.NewServer(
		weather.WithPort(9001),
		weather.WithLogLevel("quiet"),
	)

	s.GetWeather = GetNonexistentCityTestWeather
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Get("http://127.0.0.1:9001/weather?city=ZZZ")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	want := "unable to locate city \"ZZZ\"\n"
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, string(got)) {
		t.Errorf("want %q, got %q", want, string(got))
	}
}

func TestServerSeattle(t *testing.T) {
	t.Parallel()
	s := weather.NewServer(
		weather.WithPort(9002),
		weather.WithLogLevel("quiet"),
	)

	s.GetWeather = GetSeattleTestWeather
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	resp, err := http.Get("http://127.0.0.1:9002/weather?city=seattle")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	want := `{"main":"Rain","description":"Passing showers","temp":64.6,"city":"Seattle","units":"imperial"}`
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, string(got)) {
		t.Errorf("want %q, got %q", want, string(got))
	}
}
