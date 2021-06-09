package server_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
	server "weather-server"

	"github.com/google/go-cmp/cmp"
)

type Client struct {
	Base       string
	Route      string
	HTTPClient *http.Client
}

func TestServerHello(t *testing.T) {
	t.Parallel()
	//setup http server for get requests
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, "hello\n")

	}))
	defer ts.Close()

	//create new client based on struct
	var client Client

	//change default timeout value
	client.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	//set base url to test server url
	client.Base = ts.URL
	//set route to test
	client.Route = "/hello"

	url := client.Base + client.Route

	//set HTTPClient to test client to handle x509 certs w/o more setup work
	client.HTTPClient = ts.Client()

	want := "hello\n"
	resp, err := client.HTTPClient.Get(url)
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

func TestServerWeatherReport(t *testing.T) {
	t.Parallel()
	//setup http server for get requests
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		f, err := os.Open("testdata/weather_test.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		w.WriteHeader(http.StatusOK)
		io.Copy(w, f)

	}))
	defer ts.Close()

	//create new client based on struct
	var client Client

	//change default timeout value
	client.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	//set base url to test server url
	client.Base = ts.URL
	//set route to test
	client.Route = "/weatherreport"

	url := client.Base + client.Route

	//set HTTPClient to test client to handle x509 certs w/o more setup work
	client.HTTPClient = ts.Client()

	want := "hello\n"
	resp, err := client.HTTPClient.Get(url)
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

func TestPageSave(t *testing.T) {
	tmpdir := t.TempDir()
	filepath := tmpdir + "/save"
	file := tmpdir + "/save.txt"
	fmt.Println(filepath)
	testpage := server.Page{
		Title: filepath,
	}

	err := testpage.Save()
	if err != nil {
		t.Fatal(err)
	}

	_, fileerr := os.Open(file) // For read access.
	if fileerr != nil {
		log.Fatal(fileerr)
	}

}
func TestPageLoad(t *testing.T) {
	tmpdir := t.TempDir()
	filepath := tmpdir + "/load"

	fmt.Println(filepath)
	want := server.Page{
		Title: filepath,
		Body:  []byte("This is a sample Page."),
	}

	err := want.Save()
	if err != nil {
		t.Fatal(err)
	}

	got, err := server.LoadPage(want.Title)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(&want, got) {
		t.Error(cmp.Diff(want, got))
	}

}
