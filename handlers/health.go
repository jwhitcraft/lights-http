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

package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HealthHandler struct {
	Controller ControllerInterface
	Logger     *slog.Logger
	StartTime  time.Time
}

type HealthStatus struct {
	Status    string           `json:"status"`
	Timestamp time.Time        `json:"timestamp"`
	Uptime    string           `json:"uptime"`
	Checks    map[string]Check `json:"checks"`
}

type Check struct {
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r.Context())
	h.Logger.Info("Health check requested", "requestID", requestID)

	checks := make(map[string]Check)

	// Check controller and device connectivity
	if h.Controller != nil {
		devices := h.Controller.Devices()
		if len(devices) > 0 {
			checks["controller"] = Check{
				Status: "ok",
				Detail: fmt.Sprintf("%d devices connected", len(devices)),
			}
		} else {
			checks["controller"] = Check{
				Status: "ok", // Controller is working, just no devices found yet
				Detail: "Controller initialized, no devices currently connected",
			}
		}
	} else {
		checks["controller"] = Check{
			Status: "error",
			Detail: "Controller not initialized",
		}
	}

	// Overall status determination
	status := "ok"
	for _, check := range checks {
		if check.Status == "error" {
			status = "error"
			break
		}
		if check.Status == "warn" && status == "ok" {
			status = "warn"
		}
	}

	health := HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(h.StartTime).String(),
		Checks:    checks,
	}

	w.Header().Set("Content-Type", "application/json")

	// Return appropriate HTTP status
	switch status {
	case "ok":
		w.WriteHeader(http.StatusOK)
	case "warn":
		w.WriteHeader(http.StatusOK) // Warnings still return 200
	case "error":
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(health); err != nil {
		h.Logger.Error("Failed to encode health response",
			"requestID", requestID,
			"error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
