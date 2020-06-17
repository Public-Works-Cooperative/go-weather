package main

import (
	"flag"
	"fmt"
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

	start := time.Now()

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

	elapsed := time.Now().Sub(start)
	fmt.Printf("Elapsed time: %d\n", elapsed.Milliseconds())
}

func getWeatherForPlace(place string, units string, period string) (w OpenWeatherResponseOneCall, err error) {
	ll, err := getLatLngForPlace(place)
	if err != nil {
		return w, err
	}
	return getWeatherForLatLng(ll, units, period)
}

func printWeatherResult(w interface{}, place string, units string) {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	fmt.Printf("Weather for %s:\n", place)

	switch w.(type) {
	case OpenWeatherResponseCurrent:
		weath := w.(OpenWeatherResponseCurrent)
		fmt.Printf("Current: %g%s | Humidity: %d%% | %s\n",
			weath.Temp,
			unitAbbr,
			weath.Humidity,
			weath.Weather[0].Description,
		)
	case []OpenWeatherResponseHourly:
		weath := w.([]OpenWeatherResponseHourly)
		for _, h := range weath {
			t := time.Unix(h.Dt, 0)
			fmt.Printf("%-9s %2d/%2d %02d:00   %5.2f%s | Humidity: %d%% | %s\n",
				t.Weekday().String(),
				t.Month(),
				t.Day(),
				t.Hour(),
				h.Temp,
				unitAbbr,
				h.Humidity,
				h.Weather[0].Description,
			)
		}
	case []OpenWeatherResponseDaily:
		weath := w.([]OpenWeatherResponseDaily)
		for _, d := range weath {
			t := time.Unix(d.Dt, 0)
			fmt.Printf("%-9s %2d/%2d   High: %5.2f%s Low: %5.2f%s | Humidity: %d%% | %s\n",
				t.Weekday().String(),
				t.Month(),
				t.Day(),
				d.Temp.Max,
				unitAbbr,
				d.Temp.Min,
				unitAbbr,
				d.Humidity,
				d.Weather[0].Description,
			)
		}
	}
}
