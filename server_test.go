package server_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

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
