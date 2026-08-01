package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/devbrsam/fc-go-weather-by-zip-code/internal/location"
)

type locationFinderStub struct {
	city string
	err  error
}

func (s locationFinderStub) FindCity(context.Context, string) (string, error) {
	return s.city, s.err
}

type temperatureProviderStub struct {
	temperature float64
	err         error
}

func (s temperatureProviderStub) CurrentTemperature(context.Context, string) (float64, error) {
	return s.temperature, s.err
}

func TestTemperatureConversions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		celsius    float64
		fahrenheit float64
		kelvin     float64
	}{
		{name: "positive", celsius: 28.5, fahrenheit: 83.3, kelvin: 301.5},
		{name: "zero", celsius: 0, fahrenheit: 32, kelvin: 273},
		{name: "negative", celsius: -10, fahrenheit: 14, kelvin: 263},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := celsiusToFahrenheit(tt.celsius); got != tt.fahrenheit {
				t.Errorf("celsiusToFahrenheit(%v) = %v, want %v", tt.celsius, got, tt.fahrenheit)
			}
			if got := celsiusToKelvin(tt.celsius); got != tt.kelvin {
				t.Errorf("celsiusToKelvin(%v) = %v, want %v", tt.celsius, got, tt.kelvin)
			}
		})
	}
}

func TestGetWeather(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                string
		path                string
		locationFinder      LocationFinder
		temperatureProvider TemperatureProvider
		wantStatus          int
		wantBody            string
	}{
		{
			name:                "success",
			path:                "/weather/01001000",
			locationFinder:      locationFinderStub{city: "São Paulo"},
			temperatureProvider: temperatureProviderStub{temperature: 28.5},
			wantStatus:          http.StatusOK,
			wantBody:            `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}` + "\n",
		},
		{
			name:                "invalid zipcode length",
			path:                "/weather/0100100",
			locationFinder:      locationFinderStub{},
			temperatureProvider: temperatureProviderStub{},
			wantStatus:          http.StatusUnprocessableEntity,
			wantBody:            "invalid zipcode",
		},
		{
			name:                "invalid zipcode characters",
			path:                "/weather/01001-000",
			locationFinder:      locationFinderStub{},
			temperatureProvider: temperatureProviderStub{},
			wantStatus:          http.StatusUnprocessableEntity,
			wantBody:            "invalid zipcode",
		},
		{
			name:                "zipcode not found",
			path:                "/weather/99999999",
			locationFinder:      locationFinderStub{err: location.ErrZipcodeNotFound},
			temperatureProvider: temperatureProviderStub{},
			wantStatus:          http.StatusNotFound,
			wantBody:            "can not find zipcode",
		},
		{
			name:                "location service failure",
			path:                "/weather/01001000",
			locationFinder:      locationFinderStub{err: errors.New("unavailable")},
			temperatureProvider: temperatureProviderStub{},
			wantStatus:          http.StatusBadGateway,
			wantBody:            "failed to retrieve location",
		},
		{
			name:                "weather service failure",
			path:                "/weather/01001000",
			locationFinder:      locationFinderStub{city: "São Paulo"},
			temperatureProvider: temperatureProviderStub{err: errors.New("unavailable")},
			wantStatus:          http.StatusBadGateway,
			wantBody:            "failed to retrieve weather",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			handler := NewHandler(tt.locationFinder, tt.temperatureProvider)
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			handler.Routes().ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if response.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", response.Body.String(), tt.wantBody)
			}
		})
	}
}
