package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientCurrentTemperature(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/current.json" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/v1/current.json")
		}
		if got := r.URL.Query().Get("key"); got != "test-key" {
			t.Errorf("key = %q, want %q", got, "test-key")
		}
		if got := r.URL.Query().Get("q"); got != "São Paulo" {
			t.Errorf("q = %q, want %q", got, "São Paulo")
		}
		if got := r.URL.Query().Get("aqi"); got != "no" {
			t.Errorf("aqi = %q, want %q", got, "no")
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "test-key")
	temperature, err := client.CurrentTemperature(context.Background(), "São Paulo")
	if err != nil {
		t.Fatalf("CurrentTemperature returned an error: %v", err)
	}
	if temperature != 28.5 {
		t.Errorf("temperature = %v, want %v", temperature, 28.5)
	}
}

func TestClientCurrentTemperatureUnexpectedStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.Client(), server.URL, "invalid-key")
	if _, err := client.CurrentTemperature(context.Background(), "São Paulo"); err == nil {
		t.Fatal("expected an error")
	}
}
