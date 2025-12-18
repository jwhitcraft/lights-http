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

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jwhitcraft/lights-http/config"
	"github.com/jwhitcraft/lights-http/controller"
	"github.com/jwhitcraft/lights-http/handlers"
	"github.com/jwhitcraft/lights-http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	goveeController := controller.NewGoveeController(logger)

	go func() {
		err := goveeController.Start()
		if err != nil {
			logger.Error("Failed to start controller", "error", err)
		}
	}()

	defer func() {
		err := goveeController.Shutdown()
		if err != nil {
			logger.Error("Failed to shutdown controller", "error", err)
		}
		logger.Info("Controller shutdown complete")
	}()

	lightsHandler := &handlers.LightsHandler{
		Controller: goveeController.Controller,
		Logger:     logger,
	}

	healthHandler := &handlers.HealthHandler{
		Controller: goveeController.Controller,
		Logger:     logger,
		StartTime:  time.Now(),
	}

	loggingMiddleware := &middleware.LoggingMiddleware{Logger: logger}
	metricsMiddleware := &middleware.MetricsMiddleware{}

	// API server mux (with auth and metrics middleware)
	apiMux := http.NewServeMux()
	apiMux.Handle("/health", loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(healthHandler.Health))))
	apiMux.Handle("/ready", loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(healthHandler.Health))))
	apiMux.Handle("/live", loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(healthHandler.Health))))
	apiMux.Handle("/lights/on", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.TurnOn)))))
	apiMux.Handle("/lights/off", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.TurnOff)))))
	apiMux.Handle("/lights/red", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.Red)))))
	apiMux.Handle("/lights/yellow", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.Yellow)))))
	apiMux.Handle("/lights/orange", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.Orange)))))
	apiMux.Handle("/lights/dark-red", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.DarkRed)))))
	apiMux.Handle("/lights/rgb", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.RGB)))))
	apiMux.Handle("/lights/colortemp", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.ColorTemp)))))
	apiMux.Handle("/lights/brightness", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.Brightness)))))
	apiMux.Handle("/lights/status", middleware.AuthMiddleware(cfg.BearerToken)(loggingMiddleware.Middleware(metricsMiddleware.Middleware(http.HandlerFunc(lightsHandler.Status)))))

	// Metrics server mux (no auth, separate port)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())

	// Custom handler to redirect 404 and 401 to xkcd
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusRecorder{ResponseWriter: w, status: 200}
		apiMux.ServeHTTP(srw, r)
		if srw.status == 404 || srw.status == 401 {
			http.Redirect(w, r, "https://xkcd.com/random/", http.StatusFound)
		}
	})

	// Start metrics server in background
	metricsAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.MetricsPort)
	go func() {
		logger.Info("Starting metrics server", "addr", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, metricsMux); err != nil {
			logger.Error("Metrics server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Start main API server
	apiAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	logger.Info("Starting API server", "addr", apiAddr, "metrics_addr", metricsAddr)
	if err := http.ListenAndServe(apiAddr, apiHandler); err != nil {
		logger.Error("API server failed", "error", err)
		os.Exit(1)
	}
}
