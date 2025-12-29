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
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jwhitcraft/lights-http/metrics"
	govee "github.com/swrm-io/go-vee"
)

// ControllerInterface defines the methods needed for controlling lights
type ControllerInterface interface {
	Devices() []*govee.Device
}

type LightsHandler struct {
	Controller ControllerInterface
	Logger     *slog.Logger
}

// parseAndValidateJSON parses JSON from request body and validates it
func (h *LightsHandler) parseAndValidateJSON(w http.ResponseWriter, r *http.Request, v interface{}, operationName string) bool {
	requestID := getRequestID(r.Context())
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		h.Logger.Error(fmt.Sprintf("Invalid JSON in %s request", operationName),
			"requestID", requestID,
			"error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return false
	}
	return true
}

// executeLightOperation executes a light operation across all devices with proper error handling and metrics
func (h *LightsHandler) executeLightOperation(w http.ResponseWriter, r *http.Request, operationName string, successMessage string, operationFunc func(device *govee.Device) error) {
	requestID := getRequestID(r.Context())
	h.Logger.Info(fmt.Sprintf("Executing %s operation", operationName), "requestID", requestID)

	success := true
	devices := h.Controller.Devices()
	for i, device := range devices {
		if err := operationFunc(device); err != nil {
			h.Logger.Error(fmt.Sprintf("Failed to %s device", operationName),
				"device", device.DeviceID(),
				"requestID", requestID,
				"error", err)
			success = false
		}
		// Add a small delay between device operations to prevent channel blocking
		// This helps avoid "channel blocked or closed" errors when controlling multiple devices
		if i < len(devices)-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	result := "success"
	if !success {
		result = "error"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("failed to %s some lights", operationName)})
		return
	}

	metrics.LightOperationsTotal.WithLabelValues(operationName, result).Inc()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": successMessage})
}

// getRequestID safely extracts request ID from context
func getRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value("requestID").(string); ok {
		return reqID
	}
	return "unknown"
}

func (h *LightsHandler) TurnOn(w http.ResponseWriter, r *http.Request) {
	h.executeLightOperation(w, r, "turn_on", "lights turned on", func(device *govee.Device) error {
		return device.TurnOn()
	})
}

func (h *LightsHandler) TurnOff(w http.ResponseWriter, r *http.Request) {
	h.executeLightOperation(w, r, "turn_off", "lights turned off", func(device *govee.Device) error {
		return device.TurnOff()
	})
}

func (h *LightsHandler) SetColor(w http.ResponseWriter, r *http.Request, color govee.Color, colorName string) {
	h.executeLightOperation(w, r, "set_color", "lights set to "+colorName, func(device *govee.Device) error {
		return device.SetColor(color)
	})
}

func (h *LightsHandler) Red(w http.ResponseWriter, r *http.Request) {
	h.SetColor(w, r, govee.Color{R: 255, G: 0, B: 0}, "red")
}

func (h *LightsHandler) Yellow(w http.ResponseWriter, r *http.Request) {
	h.SetColor(w, r, govee.Color{R: 255, G: 255, B: 0}, "yellow")
}

func (h *LightsHandler) Orange(w http.ResponseWriter, r *http.Request) {
	h.SetColor(w, r, govee.Color{R: 139, G: 64, B: 0}, "orange")
}

func (h *LightsHandler) DarkRed(w http.ResponseWriter, r *http.Request) {
	h.SetColor(w, r, govee.Color{R: 255, G: 11, B: 0}, "dark-red")
}

func (h *LightsHandler) RGB(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r.Context())
	h.Logger.Info("Setting RGB color", "requestID", requestID)

	var req struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	}
	if !h.parseAndValidateJSON(w, r, &req, "RGB") {
		return
	}

	if req.R < 0 || req.R > 255 || req.G < 0 || req.G > 255 || req.B < 0 || req.B > 255 {
		h.Logger.Warn("Invalid RGB values",
			"requestID", requestID,
			"r", req.R, "g", req.G, "b", req.B)
		http.Error(w, "RGB values must be between 0 and 255", http.StatusBadRequest)
		return
	}

	color := govee.Color{R: uint(req.R), G: uint(req.G), B: uint(req.B)}
	h.Logger.Info("Setting RGB color",
		"requestID", requestID,
		"color", fmt.Sprintf("rgb(%d,%d,%d)", req.R, req.G, req.B))

	h.SetColor(w, r, color, "rgb")
}

func (h *LightsHandler) ColorTemp(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r.Context())
	h.Logger.Info("Setting color temperature", "requestID", requestID)

	var req struct {
		Temperature int `json:"temperature"`
	}
	if !h.parseAndValidateJSON(w, r, &req, "color temperature") {
		return
	}

	if req.Temperature < 2000 || req.Temperature > 9000 {
		h.Logger.Warn("Invalid color temperature",
			"requestID", requestID,
			"temperature", req.Temperature)
		http.Error(w, "Color temperature must be between 2000K and 9000K", http.StatusBadRequest)
		return
	}

	colorTemp := govee.NewColorKelvin(uint(req.Temperature))
	h.Logger.Info("Setting color temperature",
		"requestID", requestID,
		"temperature", fmt.Sprintf("%dK", req.Temperature))

	h.executeLightOperation(w, r, "set_color_temp", "color temperature set", func(device *govee.Device) error {
		return device.SetColorKelvin(colorTemp)
	})
}

func (h *LightsHandler) Brightness(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r.Context())
	h.Logger.Info("Setting brightness", "requestID", requestID)

	var req struct {
		Brightness int `json:"brightness"`
	}
	if !h.parseAndValidateJSON(w, r, &req, "brightness") {
		return
	}

	if req.Brightness < 0 || req.Brightness > 100 {
		h.Logger.Warn("Invalid brightness value",
			"requestID", requestID,
			"brightness", req.Brightness)
		http.Error(w, "Brightness must be between 0 and 100", http.StatusBadRequest)
		return
	}

	h.executeLightOperation(w, r, "set_brightness", "brightness set", func(device *govee.Device) error {
		return device.SetBrightness(govee.Brightness(req.Brightness))
	})
}

func (h *LightsHandler) Status(w http.ResponseWriter, r *http.Request) {
	requestID := getRequestID(r.Context())
	h.Logger.Info("Getting lights status", "requestID", requestID)

	var statuses []map[string]interface{}
	for _, device := range h.Controller.Devices() {
		err := device.RequestStatus()
		if err != nil {
			h.Logger.Error("Failed to request status", "device", device.DeviceID(), "requestID", requestID, "error", err)
			continue
		}
		color := device.Color()
		status := map[string]interface{}{
			"deviceID":   device.DeviceID(),
			"onOff":      device.Active(),
			"brightness": int(device.Brightness()),
			"color": map[string]int{
				"r": int(color.R),
				"g": int(color.G),
				"b": int(color.B),
			},
			"colortemp": device.ColorKelvin().String(),
		}
		statuses = append(statuses, status)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statuses)
}
