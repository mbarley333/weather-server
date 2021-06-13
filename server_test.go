package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	server "weather-server"

	"github.com/google/go-cmp/cmp"
)

type Client struct {
	Base       string
	Route      string
	HTTPClient *http.Client
}

// curl https://localhost:9000/weather
// curl https://localhost:9000/weather -X POST -d '{"id":"id3","main":"Sunny","description":"Clear","temp":74.6,"city":"Kyoto"}'

var wh = server.WeatherHandlers{
	Store: map[string]server.Weather{
		"id1": {
			Id:          "id1",
			Main:        "Cloudy",
			Description: "Partly cloudy",
			Temp:        74.6,
			City:        "Kaneohe",
		},
		"id2": {
			Id:          "id2",
			Main:        "Rain",
			Description: "Passing showers",
			Temp:        64.6,
			City:        "Seattle",
		},
	},
}

//TODO refactor to allow url paramters with get
//api.openweathermap.org/data/2.5/weather?q={city name}

func TestServerGet(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("GET", "/weather", nil)
	w := httptest.NewRecorder()

	wh.Get(w, r)

	want := http.StatusOK
	got := w.Code

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	wantBody := []byte(`[{"id":"id1","main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe"},{"id":"id2","main":"Rain","description":"Passing showers","temp":64.6,"city":"Seattle"}]`)
	gotBody := w.Body.Bytes()

	if !cmp.Equal(wantBody, gotBody) {
		t.Error(cmp.Diff(wantBody, gotBody))
	}

}

func TestServerPost(t *testing.T) {
	t.Parallel()

	var jsonStr = []byte(`{"id":"id3","main":"Cloudy","description":"Partly cloudy","temp":84.6,"city":"Honolulu"}`)

	r, _ := http.NewRequest("POST", "/weather", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	wh.Post(w, r)

	want := http.StatusOK
	got := w.Code

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	r2, _ := http.NewRequest("GET", "/weather", nil)
	r2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()

	wh.Get(w2, r2)

	wantBody := []byte(`[{"id":"id1","main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe"},{"id":"id2","main":"Rain","description":"Passing showers","temp":64.6,"city":"Seattle"},{"id":"id3","main":"Cloudy","description":"Partly cloudy","temp":84.6,"city":"Honolulu"}]`)
	gotBody := w2.Body.Bytes()

	if !cmp.Equal(wantBody, gotBody) {
		t.Error(cmp.Diff(wantBody, gotBody))
	}
}

func TestGetEnvironmentVariables(t *testing.T) {
	var got bool

	want := true
	if os.Getenv("CERTPATH") != "" {
		got = true
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	if os.Getenv("KEYPATH") != "" {
		got = true
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestServerGetByCity(t *testing.T) {
	t.Parallel()
	r, _ := http.NewRequest("GET", "/weather/", nil)
	w := httptest.NewRecorder()

	wh.GetByCity(w, r)

	want := http.StatusOK
	got := w.Code

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

	wantBody := []byte(`[{"id":"id1","main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe"}]`)
	gotBody := w.Body.Bytes()

	if !cmp.Equal(wantBody, gotBody) {
		t.Error(cmp.Diff(wantBody, gotBody))
	}

}
