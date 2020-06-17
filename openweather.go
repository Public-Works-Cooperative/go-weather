package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type OpenWeatherCondition struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

type OpenWeatherResponseCurrent struct {
	Dt         int64
	Sunrise    int
	Sunset     int
	Temp       float32
	Feels_like float32
	Pressure   int
	Humidity   int
	Dew_point  float32
	Uvi        float32
	Clouds     int
	Visibility int
	Wind_speed float32
	Wind_gust  float32
	Wind_deg   int
	Weather    []OpenWeatherCondition
	Rain       struct {
		_1hr float32 `json:"1hr"`
	}
}

type OpenWeatherResponseHourly struct {
	Dt         int64
	Temp       float32
	Feels_like float32
	Pressure   int
	Humidity   int
	Dew_point  float32
	Clouds     int
	Visibility int
	Wind_speed float32
	Wind_gust  float32
	Wind_deg   int
	Weather    []OpenWeatherCondition
	Rain       struct {
		_1hr float32 `json:"1hr"`
	}
}

type OpenWeatherResponseDaily struct {
	Dt      int64
	Sunrise int
	Sunset  int
	Temp    struct {
		Day   float32
		Min   float32
		Max   float32
		Night float32
		Eve   float32
		Morn  float32
	}
	Feels_like struct {
		Day   float32
		Night float32
		Eve   float32
		Morn  float32
	}
	Pressure   int
	Humidity   int
	Dew_point  float32
	Uvi        float32
	Clouds     int
	Visibility int
	Wind_speed float32
	Wind_gust  float32
	Wind_deg   int
	Weather    []OpenWeatherCondition
	Rain       float32 `json:"1hr"`
}

type OpenWeatherResponseOneCall struct {
	Current *OpenWeatherResponseCurrent
	Hourly  *[]OpenWeatherResponseHourly
	Daily   *[]OpenWeatherResponseDaily
}

const (
	WeatherPeriodCurrent  = "current"
	WeatherPeriodMinutely = "minutely"
	WeatherPeriodHourly   = "hourly"
	WeatherPeriodDaily    = "daily"
	UnitsImperial         = "imperial"
	UnitsMetric           = "metric"
)

func getWeatherForLatLng(latLng LatLng, units string, period string) (weather OpenWeatherResponseOneCall, err error) {
	// build exclude-list; always exclude minutely
	exclude := []string{WeatherPeriodMinutely}

	if period != WeatherPeriodCurrent {
		exclude = append(exclude, WeatherPeriodCurrent)
	}
	if period != WeatherPeriodHourly {
		exclude = append(exclude, WeatherPeriodHourly)
	}
	if period != WeatherPeriodDaily {
		exclude = append(exclude, WeatherPeriodDaily)
	}

	excString := strings.Join(exclude, ",")

	u := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?appid=%s&lat=%g&lon=%g&exclude=%s&units=%s",
		OpenWeatherApiKey,
		latLng.Lat,
		latLng.Lng,
		excString,
		units,
	)

	r, err := http.Get(u)
	if err != nil {
		return OpenWeatherResponseOneCall{}, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return OpenWeatherResponseOneCall{}, errors.New(fmt.Sprintf("OpenWeatherRequest Failed: %s", r.Status))
	}

	err = json.NewDecoder(r.Body).Decode(&weather)
	if err != nil {
		return OpenWeatherResponseOneCall{}, err
	}

	return weather, nil
}
