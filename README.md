# Weather Server

A Go based learning project that creates an HTTP server and queryable api endpoint.

```bash
./weather-server -port 9011 -loglevel verbose
{Clouds few clouds 81.86 Seattle}
```

## Usage
* Prior to use, install [Golang](https://golang.org/doc/install)
* Create an account on [Open Weather Map](https://home.openweathermap.org/users/sign_up) and sign up for an [API key](https://home.openweathermap.org/api_keys)
* Create an environment variable for your Open Weather Map API key: `export WEATHERAPI=YourOpenWeatherMapAPIKey`
* Clone weather repo to local machine and change to that directory
* Build weather client: go build -o weather-server ./cmd/main.go
* ./weather-server -port 9011 -loglevel verbose
* port flag can accept any port # that doesn't conflict
* loglevel flag can accept verbose or quiet
* Starts up an HTTP server on http://127.0.0.1:9011

```bash
curl http://127.0.0.1:9011/weather?city=seattle
curl http://127.0.0.1:9011/weather?city=seattle&units=imperial
{Clouds few clouds 81.86 Seattle}
```


## Goals
To learn and become more familiar with the following aspects of the Go language:
* testing
* functions and methods
* structs
* maps
* HTTP Server
* HTTP Client
* API
* concurrency
* channels
* context
* goroutines
* GUI


