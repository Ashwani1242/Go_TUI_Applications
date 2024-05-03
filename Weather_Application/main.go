package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		Temperature float64 `json:"temp_c"`
		Condition   struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch   int64   `json:"time_epoch"`
				Temperature float64 `json:"temp_c"`
				Condition   struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	city := "Lucknow"
	name := "Ashwani"
	apiKey := "40d31cxxxxxxxxxxxxxxxx11240504"  // Replace this with a valid API key.

	if len(os.Args) >= 2 {
		city = os.Args[1]
	}

	FgMagenta := color.New(color.FgHiMagenta).SprintFunc()
	FgGreen := color.New(color.FgGreen).SprintFunc()

	fmt.Printf("Hello %s\n", FgMagenta(name, "!"))

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=" + apiKey + "&q=" + city + "&days=1&")

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)

	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.ForecastDay[0].Hour

	fmt.Printf(
		"%s, %s, %s, %s\n",
		location.Name,
		location.Country,
		FgGreen(current.Temperature, "°C"),
		current.Condition.Text,
	)

	for _, hour := range hours {

		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"At %s : %.0f°C with %.0f%% rain, %s\n",
			date.Format("15:05"),
			hour.Temperature,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			color.Cyan(message)
		} else {
			color.Red(message)
		}
	}

}
