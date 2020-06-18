package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type OpenWeatherCondition struct {
	Id          int
	Main        string
	Description string
	Icon        string
}

type OpenWeatherResponseCurrent struct {
	Dt         int64
	Sunrise    int64
	Sunset     int64
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
	Snow struct {
		_1hr float32 `json:"1hr"`
	}
}

func (w OpenWeatherResponseCurrent) Output(units string) string {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	return fmt.Sprintf("Current: %g%s | Humidity: %d%% | %s\n",
		w.Temp,
		unitAbbr,
		w.Humidity,
		w.Weather[0].Description,
	)
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
	Snow struct {
		_1hr float32 `json:"1hr"`
	}
}

func (w OpenWeatherResponseHourly) Output(units string) string {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	t := time.Unix(w.Dt, 0)
	return fmt.Sprintf("%-9s %2d/%2d %02d:00   %5.2f%s | Humidity: %d%% | %s\n",
		t.Weekday().String(),
		t.Month(),
		t.Day(),
		t.Hour(),
		w.Temp,
		unitAbbr,
		w.Humidity,
		w.Weather[0].Description,
	)
}

type OpenWeatherResponseDaily struct {
	Dt      int64
	Sunrise int64
	Sunset  int64
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
	Rain       float32
	Snow       float32
}

func (w OpenWeatherResponseDaily) Output(units string) string {
	var unitAbbr string

	switch units {
	case UnitsMetric:
		unitAbbr = "C"
	case UnitsImperial:
		unitAbbr = "F"
	}

	t := time.Unix(w.Dt, 0)
	return fmt.Sprintf("%-9s %2d/%2d   High: %5.2f%s Low: %5.2f%s | Humidity: %d%% | %s\n",
		t.Weekday().String(),
		t.Month(),
		t.Day(),
		w.Temp.Max,
		unitAbbr,
		w.Temp.Min,
		unitAbbr,
		w.Humidity,
		w.Weather[0].Description,
	)
}

type OpenWeatherResponseOneCall struct {
	Current *OpenWeatherResponseCurrent
	Hourly  *[]OpenWeatherResponseHourly
	Daily   *[]OpenWeatherResponseDaily
}

func getWeatherForLatLng(latLng LatLng, units string, period string) (weather OpenWeatherResponseOneCall, err error) {
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

	r, err := httpClient.Get(u)
	if err != nil {
		return weather, err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return weather, errors.New(fmt.Sprintf("OpenWeatherRequest Failed: %s", r.Status))
	}

	err = json.NewDecoder(r.Body).Decode(&weather)

	return weather, err
}
