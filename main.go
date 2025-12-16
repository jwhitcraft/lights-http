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

	"github.com/jwhitcraft/lights-http/config"
	"github.com/jwhitcraft/lights-http/controller"
	"github.com/jwhitcraft/lights-http/handlers"
	"github.com/jwhitcraft/lights-http/middleware"
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
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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

	mux := http.NewServeMux()
	mux.Handle("/lights/on", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.TurnOn)))
	mux.Handle("/lights/off", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.TurnOff)))
	mux.Handle("/lights/red", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.Red)))
	mux.Handle("/lights/yellow", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.Yellow)))
	mux.Handle("/lights/orange", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.Orange)))
	mux.Handle("/lights/dark-red", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.DarkRed)))
	mux.Handle("/lights/rgb", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.RGB)))
	mux.Handle("/lights/brightness", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.Brightness)))
	mux.Handle("/lights/status", middleware.AuthMiddleware(cfg.BearerToken)(http.HandlerFunc(lightsHandler.Status)))

	// Custom handler to redirect 404 and 401 to xkcd
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusRecorder{ResponseWriter: w, status: 200}
		mux.ServeHTTP(srw, r)
		if srw.status == 404 || srw.status == 401 {
			http.Redirect(w, r, "https://xkcd.com/random/", http.StatusFound)
		}
	})

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	logger.Info("Starting server", "addr", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
