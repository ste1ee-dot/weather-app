package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type ForecastResponse struct {
	Meta       Meta       `json:"meta"`
	Properties Properties `json:"properties"`
}

type Meta struct {
	Version   string `json:"version"`
	License   string `json:"license"`
	Timestamp string `json:"timestamp"`
}

type Properties struct {
	Timeseries []Timeseries `json:"timeseries"`
}

type Timeseries struct {
	Time time.Time `json:"time"`
	Data Data      `json:"data"`
}

type Data struct {
	Instant     Instant     `json:"instant"`
	Next12Hours Next12Hours `json:"next_12_hours"`
}

type Instant struct {
	Details Details `json:"details"`
}

type Details struct {
	AirTemperature float64 `json:"air_temperature"`
	AirHumidity    float64 `json:"air_humidity"`
	AirPressure    float64 `json:"air_pressure"`
	WindSpeed      float64 `json:"wind_speed"`
	WindDirection  float64 `json:"wind_direction"`
}

type Next12Hours struct {
	Summary Summary `json:"summary"`
}

type Summary struct {
	SymbolCode string `json:"symbol_code"`
	Value      string `json:"value"`
}

type WeatherData struct {
	Time           string  `json:"time"`
	Temperature    float64 `json:"temperature"`
	WindSpeed      float64 `json:"wind_speed"`
	WindDirection  float64 `json:"wind_direction"`
	AirPressure    float64 `json:"air_pressure"`
	AirHumidity    float64 `json:"air_humidity"`
	WeatherSummary string  `json:"weather_summary"`
	SymbolCode     string  `json:"symbol_code"`
}

func (a *App) Log(result string) {
	fmt.Println(result)
}

// Greet returns a greeting for the given name
func (a *App) Greet(coordinates string) string {

	var err error

	coords := strings.Split(coordinates, ",")

	var myClient = &http.Client{Timeout: 10 * time.Second}

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:136.0) Gecko/20100101 Firefox/136.0",
	}

	url := "https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=" + coords[0] + "&lon=" + coords[1]

	req, err := http.NewRequest("GET", url, nil)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := myClient.Do(req)
	if err != nil {
		return fmt.Sprintf("Error getting weather info: %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading body: %s", err)
	}
	resp.Body.Close()

	var forecast ForecastResponse
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return fmt.Sprintf("Error parsing json: %s", err)
	}

	if len(forecast.Properties.Timeseries) > 0 {
		instant := forecast.Properties.Timeseries[1].Data.Instant.Details

		weatherData := WeatherData{
			Time:           forecast.Properties.Timeseries[1].Time.Format(time.RFC3339),
			Temperature:    instant.AirTemperature,
			WindSpeed:      instant.WindSpeed,
			WindDirection:  instant.WindDirection,
			AirPressure:    instant.AirPressure,
			AirHumidity:    instant.AirHumidity,
			WeatherSummary: forecast.Properties.Timeseries[1].Data.Next12Hours.Summary.Value,
			SymbolCode:     forecast.Properties.Timeseries[1].Data.Next12Hours.Summary.SymbolCode,
		}

		result, err := json.Marshal(weatherData)
		if err != nil {
			return fmt.Sprintf("Error marshaling weather data: %s", err)
		}

		return string(result)
	}
	return fmt.Sprintf("No weather data found!")
}
