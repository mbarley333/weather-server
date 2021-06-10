package server

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	testpage := &Page{
		Title: "goforit",
		Body:  []byte("I'm learning Go."),
	}

	testpage.Save()
	var dir string
	var wait time.Duration

	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/hello", Hello)
	r.HandleFunc("/weatherreport", WeatherReport)
	//not working
	r.HandleFunc("/view/", viewHandler)

	srv := &http.Server{
		Addr:              "127.0.0.1:9000",
		Handler:           r,
		IdleTimeout:       5 * time.Minute,
		ReadHeaderTimeout: time.Minute,
	}
	fmt.Println("Starting up on :9000")

	// Run our server in a goroutine so that it doesn't block.
	// https://github.com/gorilla/mux
	go func() {
		if err := srv.ListenAndServeTLS(os.Getenv("CERTPATH"), os.Getenv("KEYPATH")); err != nil {
			log.Println(err)
		}
	}()

	//create channel the understand the os.Signal comands
	c := make(chan os.Signal, 1)

	//accept CTRL+c for granceful shutdown, unblock condition
	signal.Notify(c, os.Interrupt)

	//everything below the <-c is the activity once channel unblocks
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)

}

func (p *Page) Save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0644)
}

// func LoadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return &Page{}, fmt.Errorf("unable to load page: %s", err)
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

// func ViewHandler(w http.ResponseWriter, r *http.Request) {
// 	title := r.URL.Path[len("/view/"):]
// 	p, err := LoadPage(title)
// 	if err != nil {
// 		fmt.Errorf("unable to load page: %s", err)
// 	}
// 	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
// }

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
