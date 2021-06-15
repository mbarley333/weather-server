package weather_test

import (
	"errors"
	"io"
	"log"
	"net/http"
	"testing"
	"time"
	"weather"

	"github.com/google/go-cmp/cmp"
)

// curl https://localhost:9000/weather?city=kaneohe
// curl https://localhost:9000/weather -X POST -d '{"id":"id3","main":"Sunny","description":"Clear","temp":74.6,"city":"Kyoto"}'

func GetKaneoheTestWeather(string) (weather.Weather, error) {
	return weather.Weather{
		Main:        "Cloudy",
		Description: "Partly cloudy",
		Temp:        74.6,
		City:        "Kaneohe",
	}, nil
}

func GetSeattleTestWeather(string) (weather.Weather, error) {
	return weather.Weather{
		Main:        "Rain",
		Description: "Passing showers",
		Temp:        64.6,
		City:        "Seattle",
	}, nil
}

func GetNonexistentCityTestWeather(string) (weather.Weather, error) {
	return weather.Weather{}, errors.New("no weather for you")
}

func TestServerKaneohe(t *testing.T) {
	t.Parallel()
	s := weather.NewServer(9000)
	s.GetWeather = GetKaneoheTestWeather
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:9000/weather?city=kaneohe")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	want := `{"main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe"}`
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
	s := weather.NewServer(9001)
	s.GetWeather = GetNonexistentCityTestWeather
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(50 * time.Millisecond)

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
	s := weather.NewServer(9002)
	s.GetWeather = GetSeattleTestWeather
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:9002/weather?city=seattle")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status %d", resp.StatusCode)
	}
	want := `{"main":"Rain","description":"Passing showers","temp":64.6,"city":"Seattle"}`
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, string(got)) {
		t.Errorf("want %q, got %q", want, string(got))
	}
}

