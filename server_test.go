package server_test

import (
	"encoding/json"
	"io/ioutil"
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

var wh = server.WeathersHandlers{
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

func TestServerGet(t *testing.T) {
	t.Parallel()

	sliceWeather := make([]server.Weather, len(wh.Store))

	i := 0
	for _, weather := range wh.Store {
		sliceWeather[i] = weather
		i++
	}

	want := string(`[{"id":"id1","main":"Cloudy","description":"Partly cloudy","temp":74.6,"city":"Kaneohe"},{"id":"id2","main":"Rain","description":"Passing showers","temp":64.6,"city":"Seattle"}]`)

	//marshal json for Write
	jsonBytes, err := json.Marshal(sliceWeather)
	if err != nil {
		t.Fatal(err)
	}

	//setup http server for get requests
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)

	}))

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
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

// func TestServerWeatherReport(t *testing.T) {
// 	t.Parallel()
// 	//setup http server for get requests
// 	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		f, err := os.Open("testdata/weather_test.json")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer f.Close()
// 		w.WriteHeader(http.StatusOK)
// 		io.Copy(w, f)

// 	}))
// 	defer ts.Close()

// 	//create new client based on struct
// 	var client Client

// 	//change default timeout value
// 	client.HTTPClient = &http.Client{Timeout: 10 * time.Second}
// 	//set base url to test server url
// 	client.Base = ts.URL
// 	//set route to test
// 	client.Route = "/weatherreport"

// 	url := client.Base + client.Route

// 	//set HTTPClient to test client to handle x509 certs w/o more setup work
// 	client.HTTPClient = ts.Client()

// 	want := "hello\n"
// 	resp, err := client.HTTPClient.Get(url)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	data, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	got := string(data)

// 	if !cmp.Equal(want, got) {
// 		t.Error(cmp.Diff(want, got))
// 	}

// }
