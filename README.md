# go-weather
Simple command line weather client using OpenWeather API and Google Maps Geocoding.

## Installation

1. Define `OpenWeatherApiKey` in `constants.go` 

1. Define `GoogleApiKey` in `constants.go` 

1. Enable Google Maps Geocoding API in your Maps API console.

1. `go install go-weather`

1. Ensure your Go bin directory is in your `PATH`, example `.bash_profile` line:
   - `export PATH=$PATH:$HOME/go/bin`

## Usage
```
Usage: go-weather [ -period=current|hourly|daily ] [ -units=C|F ] <location>...
  -period string
        current | hourly | daily (default "current")
  -units string
        C | F (default "C")
```

## Example
```
$ go-weather -period=daily -units=f 80209 80919
Weather for 80209:
Current: 63.37F | Humidity: 36% | few clouds
Weather for 80919:
Current: 58.48F | Humidity: 33% | clear sky
```
