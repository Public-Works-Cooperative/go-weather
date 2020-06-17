package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type LatLng struct {
	Lat float64
	Lng float64
}

type GoogleGeocodeResult struct {
	Geometry struct {
		Location LatLng
	}
}

func (g GoogleGeocodeResult) ToLatLng() LatLng {
	return g.Geometry.Location
}

type GoogleGeocodeResponse struct {
	Status  string
	Results []GoogleGeocodeResult
}

func getLatLngForPlace(place string) (latLng LatLng, err error) {
	escPlace := url.QueryEscape(place)
	u := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&address=%s",
		GoogleApiKey,
		escPlace,
	)

	r, err := http.Get(u)
	if err != nil {
		return LatLng{}, err
	}
	defer r.Body.Close()

	var geocode GoogleGeocodeResponse

	err = json.NewDecoder(r.Body).Decode(&geocode)
	if err != nil {
		return LatLng{}, err
	}

	if geocode.Status != "OK" || len(geocode.Results) < 1 {
		return LatLng{}, err
	}

	return geocode.Results[0].ToLatLng(), nil
}

func printWeatherResult(w interface{}, location string, units string) {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	fmt.Printf("Weather for %s:\n", location)

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

func main() {
	ll, err := getLatLngForPlace("80919")
	if err != nil {
		log.Fatal(err)
		return
	}
	w, err := getWeatherForLatLng(ll, UnitsMetric, WeatherPeriodHourly)
	if err != nil {
		log.Fatal(err)
	}
	printWeatherResult(*w.Hourly, "80919", UnitsMetric)
}
