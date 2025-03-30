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

type LocationResponse struct {
	Address Address `json:"address"`
}
type Address struct {
	City        string `json:"city"`
	CountryCode string `json:"country_code"`
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
	Instant    Instant    `json:"instant"`
	Next1Hours Next1Hours `json:"next_1_hours"`
	Next6Hours Next6Hours `json:"next_6_hours"`
}

type Instant struct {
	Details Details `json:"details"`
}

type Details struct {
	AirTemperature float64 `json:"air_temperature"`
	AirHumidity    float64 `json:"relative_humidity"`
	AirPressure    float64 `json:"air_pressure_at_sea_level"`
	WindSpeed      float64 `json:"wind_speed"`
	WindDirection  float64 `json:"wind_from_direction"`
}

type Next1Hours struct {
	Summary Summary `json:"summary"`
}

type Next6Hours struct {
	Summary Summary `json:"summary"`
}

type Summary struct {
	SymbolCode string `json:"symbol_code"`
}

type WeatherData struct {
	Time               string  `json:"time"`
	Temperature        float64 `json:"temperature"`
	WindSpeed          float64 `json:"wind_speed"`
	WindDirection      float64 `json:"wind_direction"`
	AirPressure        float64 `json:"air_pressure"`
	AirHumidity        float64 `json:"air_humidity"`
	SymbolCode         string  `json:"symbol_code"`
	SymbolCodeNice     string  `json:"symbol_code_nice"`
	AddressCity        string  `json:"city"`
	AddressCountryCode string  `json:"country_code"`
	CurrentDay         string  `json:"week_day"`
	Date               string  `json:"date"`
	FirstDay           string  `json:"first_day"`
	SecondDay          string  `json:"second_day"`
	SecondTemp         float64 `json:"second_temp"`
	SecondSymbol       string  `json:"second_symbol"`
	ThirdDay           string  `json:"third_day"`
	ThirdTemp          float64 `json:"third_temp"`
	ThirdSymbol        string  `json:"third_symbol"`
	FourthDay          string  `json:"fourth_day"`
	FourthTemp         float64 `json:"fourth_temp"`
	FourthSymbol       string  `json:"fourth_symbol"`
}

func (a *App) Log(toLog string) {
	fmt.Println("---------")
	fmt.Println(toLog)
	fmt.Println("---------")
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

	urlLocation := "https://nominatim.openstreetmap.org/reverse?lat=" + coords[0] + "&lon=" + coords[1] + "&format=jsonv2"

	reqLocation, err := http.NewRequest("GET", urlLocation, nil)

	req, err := http.NewRequest("GET", url, nil)

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	respLoc, err := myClient.Do(reqLocation)
	if err != nil {
		return fmt.Sprintf("Error getting location info: %s", err)
	}

	bodyLoc, err := io.ReadAll(respLoc.Body)
	if err != nil {
		return fmt.Sprintf("Error reading location body: %s", err)
	}
	respLoc.Body.Close()

	var locResponse LocationResponse
	err = json.Unmarshal(bodyLoc, &locResponse)
	if err != nil {
		return fmt.Sprintf("Error parsing location json: %s", err)
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

	t, err := time.Parse(time.RFC3339, forecast.Properties.Timeseries[1].Time.Format(time.RFC3339))
	if err != nil {
		fmt.Println("Error parsing date:", err)
	}

	dayOfWeek := t.Weekday()

	firstDay := dayOfWeek

	secondTime := t.AddDate(0, 0, 1)
	secondDay := secondTime.Weekday()

	thirdTime := secondTime.AddDate(0, 0, 1)
	thirdDay := thirdTime.Weekday()

	fourthTime := thirdTime.AddDate(0, 0, 1)
	fourthDay := fourthTime.Weekday()

	formattedDate := t.Format("02 Jan 2006")

	if len(forecast.Properties.Timeseries) > 0 {
		instant := forecast.Properties.Timeseries[0].Data.Instant.Details

		weatherData := WeatherData{
			Time:               forecast.Properties.Timeseries[0].Time.Format(time.RFC3339),
			Temperature:        instant.AirTemperature,
			WindSpeed:          instant.WindSpeed,
			WindDirection:      instant.WindDirection,
			AirPressure:        instant.AirPressure,
			AirHumidity:        instant.AirHumidity,
			SymbolCode:         forecast.Properties.Timeseries[0].Data.Next1Hours.Summary.SymbolCode,
			SymbolCodeNice:     strings.Split(forecast.Properties.Timeseries[0].Data.Next1Hours.Summary.SymbolCode, "_")[0],
			AddressCity:        locResponse.Address.City,
			AddressCountryCode: strings.ToUpper(locResponse.Address.CountryCode),
			CurrentDay:         dayOfWeek.String(),
			Date:               formattedDate,
			FirstDay:           firstDay.String()[:3],
			SecondDay:          secondDay.String()[:3],
			SecondTemp:         forecast.Properties.Timeseries[24].Data.Instant.Details.AirTemperature,
			SecondSymbol:       forecast.Properties.Timeseries[24].Data.Next1Hours.Summary.SymbolCode,
			ThirdDay:           thirdDay.String()[:3],
			ThirdTemp:          forecast.Properties.Timeseries[48].Data.Instant.Details.AirTemperature,
			ThirdSymbol:        forecast.Properties.Timeseries[48].Data.Next1Hours.Summary.SymbolCode,
			FourthDay:          fourthDay.String()[:3],
			FourthTemp:         forecast.Properties.Timeseries[72].Data.Instant.Details.AirTemperature,
			FourthSymbol:       forecast.Properties.Timeseries[66].Data.Next6Hours.Summary.SymbolCode,
		}

		result, err := json.Marshal(weatherData)
		if err != nil {
			return fmt.Sprintf("Error marshaling weather data: %s", err)
		}
		return string(result)
	}
	return fmt.Sprintf("No weather data found!")
}
