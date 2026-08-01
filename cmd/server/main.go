package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/devbrsam/fc-go-weather-by-zip-code/internal/api"
	"github.com/devbrsam/fc-go-weather-by-zip-code/internal/location"
	"github.com/devbrsam/fc-go-weather-by-zip-code/internal/weather"
)

func main() {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	handler := api.NewHandler(
		location.NewViaCEPClient(client, "https://viacep.com.br"),
		weather.NewClient(client, "https://api.weatherapi.com", apiKey),
	)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
