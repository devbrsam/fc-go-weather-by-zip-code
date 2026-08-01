package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(httpClient *http.Client, baseURL, apiKey string) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
	}
}

func (c *Client) CurrentTemperature(ctx context.Context, city string) (float64, error) {
	endpoint, err := url.Parse(c.baseURL + "/v1/current.json")
	if err != nil {
		return 0, fmt.Errorf("parse WeatherAPI URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("key", c.apiKey)
	query.Set("q", city)
	query.Set("aqi", "no")
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("create WeatherAPI request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("WeatherAPI returned status %d", resp.StatusCode)
	}

	var result struct {
		Current struct {
			Temperature float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("decode WeatherAPI response: %w", err)
	}

	return result.Current.Temperature, nil
}
