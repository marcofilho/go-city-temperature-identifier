package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"unicode"

	"github.com/spf13/viper"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type conf struct {
	API_KEY string `mapstructure:"API_KEY"`
}

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
	http.HandleFunc("/cep", cepHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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

func loadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, err
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" {
		writeHTMLError(w, http.StatusBadRequest, "CEP is required! Example: /cep?cep=22222-222")
		return
	}

	if len(cep) < 9 {
		writeHTMLError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	config, err := loadConfig(".")
	if err != nil {
		writeHTMLError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	apiKey := config.API_KEY

	c, err := getCep(cep)
	if err != nil || c == (Cep{}) {
		writeHTMLError(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	type viacepError struct {
		Erro bool `json:"erro"`
	}
	var errCheck viacepError

	cepJson, _ := json.Marshal(c)
	json.Unmarshal(cepJson, &errCheck)
	if errCheck.Erro {
		writeHTMLError(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	location := removeAccents(c.Localidade)

	weather, err := getWeather(location, apiKey)
	if err != nil {
		writeHTMLError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response := map[string]interface{}{
		"city_name":    weather.Location.Name,
		"city_region":  weather.Location.Region,
		"city_country": weather.Location.Country,
		"temp_C":       weather.Current.TempC,
		"temp_F":       convertToFahrenheit(weather.Current.TempC),
		"temp_K":       convertToKelvin(weather.Current.TempC),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func writeHTMLError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	fmt.Fprintf(w, "<html><body><h2>Error: %s</h2></body></html>", message)
}
