package here

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWeatherService(t *testing.T) {
	var httpClient = &http.Client{
		Timeout: time.Second * 15,
	}
	client := NewWeatherClient(httpClient)
	assert.NotNil(t, client)
}

func TestWeatherService_SevereWeatherAlerts(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(readJSONFile("resources/severe_weather_alerts.json"))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client := NewWeatherClient(httpClient)
	weatherAlerts := SevereWeatherAlertsParams{
		Name:   "Boston",
		APIKey: "apiKey",
	}
	severeWeatherAlerts, _, err := client.Weather.SevereWeatherAlerts(&weatherAlerts)
	assert.NotNil(t, severeWeatherAlerts)
	assert.Nil(t, err)
}

func TestWeatherService_WeatherConditions(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(readJSONFile("resources/weather_conditions_by_location.json"))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	client := NewWeatherClient(httpClient)
	weatherConditionsParams := WeatherConditionsParams{
		Latitude:       52.516,
		Longitude:      13.389,
		OneObservation: true,
		APIKey:         "apiKey",
	}
	weatherConditions, _, err := client.Weather.WeatherConditions(&weatherConditionsParams)
	assert.NotNil(t, weatherConditions)
	assert.Nil(t, err)
}
