package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
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

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", Hello)
	r.HandleFunc("/weatherreport", WeatherReport)

	srv := &http.Server{
		Addr:              "127.0.0.1:9000",
		Handler:           r,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}
	fmt.Println("Starting up on :9000")

	log.Fatal(srv.ListenAndServeTLS(os.Getenv("CERTPATH"), os.Getenv("KEYPATH")))

}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return &Page{}, fmt.Errorf("unable to load page: %s", err)
	}
	return &Page{Title: title, Body: body}, err
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := LoadPage(title)
	if err != nil {
		fmt.Errorf("unable to load page: %s", err)
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
