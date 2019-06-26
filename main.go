package main

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
	"os/user"
	"time"
)

type CurrentWeather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure float64     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  float64     `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

var (
	city = kingpin.Arg("city", "Weather for the city").String()
	apikey = kingpin.Flag("api-key", "The apikey for OpenWeatherMap").String()
	current = kingpin.Flag("current", "Shows the current weather").Short('c').Bool()
	save = kingpin.Flag("save", "Saves the api key").Bool()
)

func geticon(pic string) string {
	switch pic {
	case "01d":
		return "☼"
	case "02d":
		return "🌤"
	case "03d":
		return "🌥"
	case "04d":
		return "☁"
	case "09d":
		return "🌧"
	case "10d":
		return "🌦"
	case "11d":
		return "🌩"
	case "13d":
		return "🌨"
	case "50d":
		return "🌫"
	case "01n":
		return "🌑"
	case "02n":
		return "🌤"
	case "03n":
		return "🌥"
	case "04n":
		return "☁"
	case "09n":
		return "🌧"
	case "10n":
		return "🌦"
	case "11n":
		return "🌩"
	case "13n":
		return "🌨"
	case "50n":
		return "🌫"
	}
	return ""
}

func main() {
	kingpin.Parse()

	t := time.Now()

	if *save {
		if *apikey != "" {
			key := []byte(*apikey)
			usr, _ := user.Current()
			err := ioutil.WriteFile(fmt.Sprintf("%s/.config/wa", usr.HomeDir), key, 0644)
			if err != nil {
				fmt.Println("Couldn't save api-key")
				os.Exit(1)
			}
			os.Exit(0)
		}
	}

	if *city == "" {
		fmt.Println("Please enter a City")
		os.Exit(1)
	}

	if *apikey == "" {
		usr, _ := user.Current()
		key, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/wa", usr.HomeDir))
		if err == nil {
			*apikey = string(key)
		} else {
			fmt.Println("Use an api-key!")
			os.Exit(1)
		}
	}

	if *current {
		currentweather := CurrentWeather{}
		currentweatherUrl := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", *city, *apikey)
		resp, err := req.Get(currentweatherUrl)
		if err != nil {
			panic(err)
		}

		if resp.Response().StatusCode == 404 {
			fmt.Println("City not found")
			os.Exit(1)
		}

		json.Unmarshal(resp.Bytes(), &currentweather)

		/*
			fmt.Println("temp:", currentweather.Main.Temp)
			fmt.Println("max_temp:", currentweather.Main.TempMax)
			fmt.Println("min_temp:", currentweather.Main.TempMin)
			fmt.Println("weather:", geticon(currentweather.Weather[0].Icon))
		*/

		data := [][]string{
			{t.Format("2006-01-02 15:04"), fmt.Sprintf("%v", currentweather.Main.Temp), fmt.Sprintf("%v", currentweather.Main.TempMax), fmt.Sprintf("%v", currentweather.Main.TempMin), fmt.Sprintf("%v - %v", geticon(currentweather.Weather[0].Icon), currentweather.Weather[0].Description)},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{currentweather.Name, "Current Temp", "Max Temp", "Min Temp", "Weather"})
		table.AppendBulk(data)
		table.Render()
	}
}
