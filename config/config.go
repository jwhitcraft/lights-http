// Copyright 2025 Jon Whitcraft
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config holds server and auth configuration
type Config struct {
	Host        string
	Port        string
	MetricsPort string
	BearerToken string
}

// Load loads configuration from environment variables and .env (if not production)
func Load() (*Config, error) {
	if os.Getenv("GO_ENV") != "production" {
		_ = godotenv.Load()
	}

	host := os.Getenv("HOSTNAME")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}
	token := os.Getenv("BEARER_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BEARER_TOKEN is required. Please set it in your environment or .env file")
	}

	return &Config{
		Host:        host,
		Port:        port,
		MetricsPort: metricsPort,
		BearerToken: token,
	}, nil
}
