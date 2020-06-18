package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	WeatherPeriodCurrent  = "current"
	WeatherPeriodMinutely = "minutely"
	WeatherPeriodHourly   = "hourly"
	WeatherPeriodDaily    = "daily"
	UnitsImperial         = "imperial"
	UnitsMetric           = "metric"
)

func exitInvalidArguments() {
	println("\nUsage: go-weather [ -period=current|hourly|daily ] [ -units=C|F ] <location>\n")
	flag.Usage()
	println()
	os.Exit(2)
}

func main() {
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

	// start := time.Now()

	for _, p := range places {
		w, err := getWeatherForPlace(p, un, *period)
		if err != nil {
			panic(err)
		}

		switch *period {
		case WeatherPeriodCurrent:
			printWeatherResult(*w.Current, p, un)
		case WeatherPeriodHourly:
			printWeatherResult(*w.Hourly, p, un)
		case WeatherPeriodDaily:
			printWeatherResult(*w.Daily, p, un)
		}
	}

	// elapsed := time.Now().Sub(start)
	// fmt.Printf("Elapsed time: %d\n", elapsed.Milliseconds())
}

func getWeatherForPlace(place string, units string, period string) (w OpenWeatherResponseOneCall, err error) {
	ll, err := getLatLngForPlace(place)
	if err != nil {
		return w, err
	}
	return getWeatherForLatLng(ll, units, period)
}

func printWeatherResult(w interface{}, place string, units string) {
	fmt.Printf("Weather for %s:\n", place)

	switch w.(type) {
	case OpenWeatherResponseCurrent:
		weath := w.(OpenWeatherResponseCurrent)
		fmt.Print(weath.Output(units))
	case []OpenWeatherResponseHourly:
		weath := w.([]OpenWeatherResponseHourly)
		for _, h := range weath {
			fmt.Print(h.Output(units))
		}
	case []OpenWeatherResponseDaily:
		weath := w.([]OpenWeatherResponseDaily)
		for _, d := range weath {
			fmt.Print(d.Output(units))
		}
	}
}
