package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		wantErr  bool
		expected *Config
	}{
		{
			name: "valid config with defaults",
			env: map[string]string{
				"BEARER_TOKEN": "test-token",
			},
			wantErr: false,
			expected: &Config{
				Host:        "0.0.0.0",
				Port:        "8080",
				BearerToken: "test-token",
			},
		},
		{
			name: "valid config with custom values",
			env: map[string]string{
				"HOSTNAME":     "127.0.0.1",
				"PORT":         "3000",
				"BEARER_TOKEN": "custom-token",
			},
			wantErr: false,
			expected: &Config{
				Host:        "127.0.0.1",
				Port:        "3000",
				BearerToken: "custom-token",
			},
		},
		{
			name:    "missing bearer token",
			env:     map[string]string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env
			os.Unsetenv("HOSTNAME")
			os.Unsetenv("PORT")
			os.Unsetenv("BEARER_TOKEN")
			os.Unsetenv("GO_ENV")

			// Set test env
			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && cfg != nil {
				if cfg.Host != tt.expected.Host || cfg.Port != tt.expected.Port || cfg.BearerToken != tt.expected.BearerToken {
					t.Errorf("Load() = %v, want %v", cfg, tt.expected)
				}
			}
		})
	}
}
