/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package info

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	countryName string
)

type Weather struct{
	Location struct {
		Name string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64 `json:"time_epoch"`
				TempC float64 `json:"temp_c"`
				TempF float64 `json:"temp_f"`
				Condition struct {
						Text string `json:"text"`
					} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`

			} `json:"hour"`

		} `json:"forecastday"`
	} `json:"forecast"`
}

// weatherCmd represents the weather command
var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "this command fetch's current weather of the city",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		res, err := http.Get("https://api.weatherapi.com/v1/forecast.json?key=a713fa03110942debe931017241606&q=London&days=1&aqi=no&alerts=no")

		if err != nil {
			panic(err)
		}

		if(res.StatusCode != 200){
			panic("Weather api is not available")
		}

		body, err := io.ReadAll(res.Body);

		if err != nil {
			panic(err)
		}

		var weather Weather
		err = json.Unmarshal(body, &weather)

		if( err != nil){
			panic(err)
		}

		location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

		fmt.Printf("%s, %s: %.0fC, %s \n",location.Name, location.Country, current.TempC, current.Condition.Text)

		for _, hour := range hours{
			date := time.Unix(hour.TimeEpoch,0)

			if(date.Before(time.Now())){
				continue
			}

			message := fmt.Sprintf(
				"%s - %.0fC, %.0fC, %s \n",
				date.Format("15:04"),
				hour.TempC,
				hour.ChanceOfRain,
				hour.Condition.Text,
			)

			if(hour.ChanceOfRain< 40){
				fmt.Print(message)
			}else { 
				color.Red(message)
			}
		}
	},
}

func init() {
	
	weatherCmd.Flags().StringVarP(&countryName, "city", "c", "", "set country to fetch city specific weather info")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// weatherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// weatherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
