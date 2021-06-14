package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
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

func TestServerGetByCity(t *testing.T) {
	t.Parallel()

	type testCase struct {
		url  string
		want []byte
	}
	tcs := []testCase{
		{
			url:  "/weather?city=kaneohe",
			want: []byte(`[{"id":"id1","main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe"}]`),
		},
		{
			url:  "/weather?city=zzz",
			want: []byte("unable to locate city"),
		},
	}
	for _, tc := range tcs {

		r, _ := http.NewRequest("GET", tc.url, nil)
		w := httptest.NewRecorder()

		wh.GetByCity(w, r)

		want := http.StatusOK
		got := w.Code

		if !cmp.Equal(want, got) {
			t.Error(cmp.Diff(want, got))
		}

		wantBody := tc.want
		gotBody := w.Body.Bytes()

		if !cmp.Equal(wantBody, gotBody) {
			t.Error(cmp.Diff(wantBody, gotBody))
		}

	}

}

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

func TestNewId(t *testing.T) {
	t.Parallel()

	want := true

	sliceid := make([]string, 0, 10)

	for i := 0; i < 10; i++ {
		sliceid = append(sliceid, server.NewId())
	}

	sort.Strings(sliceid)

	got := true
	for j := 1; j < 10; j++ {
		if sliceid[j] == sliceid[j-1] {
			got = false
		}
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
