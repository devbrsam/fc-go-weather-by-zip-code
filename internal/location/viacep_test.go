package location

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViaCEPClientFindCity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		status     int
		body       string
		wantCity   string
		wantErr    error
		wantErrSet bool
	}{
		{
			name:     "found",
			status:   http.StatusOK,
			body:     `{"cep":"01001-000","localidade":"São Paulo"}`,
			wantCity: "São Paulo",
		},
		{
			name:    "not found",
			status:  http.StatusOK,
			body:    `{"erro":"true"}`,
			wantErr: ErrZipcodeNotFound,
		},
		{
			name:    "empty city",
			status:  http.StatusOK,
			body:    `{"localidade":""}`,
			wantErr: ErrZipcodeNotFound,
		},
		{
			name:       "unexpected status",
			status:     http.StatusServiceUnavailable,
			body:       `{}`,
			wantErrSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/ws/01001000/json/" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/ws/01001000/json/")
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewViaCEPClient(server.Client(), server.URL)
			city, err := client.FindCity(context.Background(), "01001000")

			if city != tt.wantCity {
				t.Errorf("city = %q, want %q", city, tt.wantCity)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErrSet && err == nil {
				t.Error("expected an error")
			}
			if tt.wantErr == nil && !tt.wantErrSet && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
