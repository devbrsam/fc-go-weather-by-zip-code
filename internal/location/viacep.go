package location

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var ErrZipcodeNotFound = errors.New("zipcode not found")

type ViaCEPClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewViaCEPClient(httpClient *http.Client, baseURL string) *ViaCEPClient {
	return &ViaCEPClient{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (c *ViaCEPClient) FindCity(ctx context.Context, zipcode string) (string, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/ws/%s/json/", c.baseURL, zipcode),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("create ViaCEP request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return "", ErrZipcodeNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ViaCEP returned status %d", resp.StatusCode)
	}

	var result struct {
		City  string `json:"localidade"`
		Error any    `json:"erro"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode ViaCEP response: %w", err)
	}
	if result.Error != nil || strings.TrimSpace(result.City) == "" {
		return "", ErrZipcodeNotFound
	}

	return result.City, nil
}
