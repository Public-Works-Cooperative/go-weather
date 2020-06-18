package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	WeatherPeriodCurrent  = "current"
	WeatherPeriodMinutely = "minutely"
	WeatherPeriodHourly   = "hourly"
	WeatherPeriodDaily    = "daily"
	UnitsImperial         = "imperial"
	UnitsMetric           = "metric"
)

var httpClient http.Client

func exitInvalidArguments() {
	println("\nUsage: go-weather [ -period=current|hourly|daily ] [ -units=C|F ] <location>...\n")
	flag.Usage()
	println()
	os.Exit(2)
}

func main() {
	httpClient = http.Client{
		Timeout: time.Second * 10,
	}

	units := flag.String("units", "C", "C | F")
	period := flag.String("period", "current", "current | hourly | daily")
	flag.Parse()

	places := flag.Args()

	if len(places) < 1 {
		exitInvalidArguments()
	}

	var un string
	if strings.ToUpper(*units) == "C" {
		un = UnitsMetric
	} else if strings.ToUpper(*units) == "F" {
		un = UnitsImperial
	} else {
		exitInvalidArguments()
	}

	if *period != WeatherPeriodCurrent &&
		*period != WeatherPeriodHourly &&
		*period != WeatherPeriodDaily {
		exitInvalidArguments()
	}

	chs := make([]chan OpenWeatherResponseOneCall, len(places))
	errChs := make([]chan error, len(places))

	start := time.Now()

	for i, place := range places {
		chs[i] = make(chan OpenWeatherResponseOneCall, 1)
		errChs[i] = make(chan error, 1)
		go concurrentGetWeatherForPlace(place, un, *period, chs[i], errChs[i])
	}

	for i, ch := range chs {
		w := <-ch
		err := <-errChs[i]
		if err != nil {
			log.Fatal(err)
		} else {
			switch *period {
			case WeatherPeriodCurrent:
				printWeatherResult(*w.Current, places[i], un)
			case WeatherPeriodHourly:
				printWeatherResult(*w.Hourly, places[i], un)
			case WeatherPeriodDaily:
				printWeatherResult(*w.Daily, places[i], un)
			}
		}
	}

	elasped := time.Now().Sub(start)
	fmt.Printf("Elasped time: %d\n", elasped.Milliseconds())
}

func getWeatherForPlace(place string, units string, period string) (w OpenWeatherResponseOneCall, err error) {
	ll, err := getLatLngForPlace(place)
	if err != nil {
		return w, err
	}
	return getWeatherForLatLng(ll, units, period)
}

func concurrentGetWeatherForPlace(place string, units string, period string, wCh chan OpenWeatherResponseOneCall, errCh chan error) {
	w, err := getWeatherForPlace(place, units, period)
	wCh <- w
	errCh <- err
}

func printWeatherResult(w interface{}, place string, units string) {
	fmt.Printf("Weather for %s:\n", place)

	switch w.(type) {
	case OpenWeatherResponseCurrent:
		fmt.Print(w.(OpenWeatherResponseCurrent).Output(units))
	case []OpenWeatherResponseHourly:
		for _, h := range w.([]OpenWeatherResponseHourly) {
			fmt.Print(h.Output(units))
		}
	case []OpenWeatherResponseDaily:
		for _, h := range w.([]OpenWeatherResponseDaily) {
			fmt.Print(h.Output(units))
		}
	}
}
