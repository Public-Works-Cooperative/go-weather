package main

import (
	"encoding/json"
	"fmt"
	"net/url"
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

	r, err := httpClient.Get(u)
	if err != nil {
		return latLng, err
	}
	defer r.Body.Close()

	var geocode GoogleGeocodeResponse

	err = json.NewDecoder(r.Body).Decode(&geocode)
	if err != nil {
		return latLng, err
	}

	if geocode.Status != "OK" || len(geocode.Results) < 1 {
		return latLng, err
	}

	return geocode.Results[0].ToLatLng(), nil
}
