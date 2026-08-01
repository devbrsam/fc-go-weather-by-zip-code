package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/devbrsam/fc-go-weather-by-zip-code/internal/location"
)

var zipcodePattern = regexp.MustCompile(`^\d{8}$`)

type LocationFinder interface {
	FindCity(ctx context.Context, zipcode string) (string, error)
}

type TemperatureProvider interface {
	CurrentTemperature(ctx context.Context, city string) (float64, error)
}

type Handler struct {
	locationFinder      LocationFinder
	temperatureProvider TemperatureProvider
}

type temperatureResponse struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func NewHandler(locationFinder LocationFinder, temperatureProvider TemperatureProvider) *Handler {
	return &Handler{
		locationFinder:      locationFinder,
		temperatureProvider: temperatureProvider,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{zipcode}", h.getWeather)
	return mux
}

func (h *Handler) getWeather(w http.ResponseWriter, r *http.Request) {
	zipcode := r.PathValue("zipcode")
	if !zipcodePattern.MatchString(zipcode) {
		writeError(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	city, err := h.locationFinder.FindCity(r.Context(), zipcode)
	if errors.Is(err, location.ErrZipcodeNotFound) {
		writeError(w, http.StatusNotFound, "can not find zipcode")
		return
	}
	if err != nil {
		log.Printf("find city: %v", err)
		writeError(w, http.StatusBadGateway, "failed to retrieve location")
		return
	}

	celsius, err := h.temperatureProvider.CurrentTemperature(r.Context(), city)
	if err != nil {
		log.Printf("retrieve temperature: %v", err)
		writeError(w, http.StatusBadGateway, "failed to retrieve weather")
		return
	}

	response := temperatureResponse{
		Celsius:    celsius,
		Fahrenheit: celsiusToFahrenheit(celsius),
		Kelvin:     celsiusToKelvin(celsius),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func celsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(message))
}
