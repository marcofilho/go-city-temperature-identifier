package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"unicode"

	"github.com/marcofilho/go-city-temperature-identifier/configs"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Current struct {
	LastUpdatedEpoch int       `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Weather struct {
	Location   Location `json:"location"`
	Current    Current  `json:"current"`
	WindMph    float64  `json:"wind_mph"`
	WindKph    float64  `json:"wind_kph"`
	WindDegree int      `json:"wind_degree"`
	WindDir    string   `json:"wind_dir"`
	PressureMb int      `json:"pressure_mb"`
	PressureIn float64  `json:"pressure_in"`
	PrecipMm   float64  `json:"precip_mm"`
	PrecipIn   float64  `json:"precip_in"`
	Humidity   int      `json:"humidity"`
	Cloud      int      `json:"cloud"`
	FeelslikeC float64  `json:"feelslike_c"`
	FeelslikeF float64  `json:"feelslike_f"`
	VisKm      float64  `json:"vis_km"`
	VisMiles   float64  `json:"vis_miles"`
	Uv         float64  `json:"uv"`
	GustMph    float64  `json:"gust_mph"`
	GustKph    float64  `json:"gust_kph"`
}

type Response struct {
	City       string  `json:"city"`
	Temp_C     float64 `json:"temp"`
	Temp_F     float64 `json:"temp_f"`
	Temp_K     float64 `json:"temp_k"`
	StatusCode int     `json:"status_code"`
}

type Cep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	cep := flag.String("cep", "", "CEP to be searched")

	flag.Parse()

	if *cep == "" {
		fmt.Println("Error: CEP is required.")
		fmt.Println("Example: -cep=22222-222")
		os.Exit(1)
	}

	if len(*cep) < 9 {
		fmt.Println("Código HTTP: 422")
		fmt.Println("Mensagem: invalid zipcode")
		return
	}

	configs, err := configs.LoadConfig(".")
	apiKey := configs.API_KEY
	if err != nil {
		fmt.Println("Código HTTP: 500")
		fmt.Println(`Mensagem: internal server error`)
		return
	}

	c, err := getCep(*cep)
	if err != nil || c == (Cep{}) {
		fmt.Println("Código HTTP: 404")
		fmt.Println(`Mensagem: can not find zipcode`)
		return
	}

	type viacepError struct {
		Erro bool `json:"erro"`
	}
	var errCheck viacepError

	cepJson, _ := json.Marshal(c)
	json.Unmarshal(cepJson, &errCheck)
	if errCheck.Erro {
		fmt.Println("Código HTTP: 404")
		fmt.Println(`Mensagem: can not find zipcode`)
		return
	}

	location := removeAccents(c.Localidade)

	weather, err := getWeather(location, apiKey)
	if err != nil {
		fmt.Println("Código HTTP: 500")
		fmt.Println(`Mensagem: internal server error`)
		return
	}

	response := map[string]float64{
		"temp_C": weather.Current.TempC,
		"temp_F": convertToFahrenheit(weather.Current.TempC),
		"temp_K": convertToKelvin(weather.Current.TempC),
	}
	respJson, _ := json.Marshal(response)
	fmt.Println("Código HTTP: 200")
	fmt.Println(`City Name: ` + weather.Location.Name)
	fmt.Println(`City Region: ` + weather.Location.Region)
	fmt.Println(`City Country: ` + weather.Location.Country)
	fmt.Printf("Response Body: %s\n", respJson)
}

func getCep(cep string) (Cep, error) {
	req, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return Cep{}, err
	}

	var c Cep
	err = json.NewDecoder(req.Body).Decode(&c)
	if err != nil {
		return Cep{}, err
	}
	defer req.Body.Close()

	return c, nil
}

func getWeather(city, apiKey string) (Weather, error) {
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedCity)

	req, err := http.Get(url)
	if err != nil {
		return Weather{}, err
	}
	defer req.Body.Close()

	bodyBytes, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var w Weather
	err = json.NewDecoder(req.Body).Decode(&w)
	if err != nil {
		return Weather{}, err
	}
	return w, nil
}

func convertToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func convertToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return output
}
