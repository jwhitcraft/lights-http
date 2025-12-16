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
	"log/slog"
	"net/http"

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

func (h *LightsHandler) forEachDevice(fn func(device *govee.Device)) {
	for _, device := range h.Controller.Devices() {
		fn(device)
	}
}

func (h *LightsHandler) TurnOn(w http.ResponseWriter, r *http.Request) {
	h.forEachDevice(func(device *govee.Device) {
		err := device.TurnOn()
		if err != nil {
			h.Logger.Error("Failed to turn on device", "device", device.DeviceID(), "error", err)
		}
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "lights turned on"})
}

func (h *LightsHandler) TurnOff(w http.ResponseWriter, r *http.Request) {
	h.forEachDevice(func(device *govee.Device) {
		err := device.TurnOff()
		if err != nil {
			h.Logger.Error("Failed to turn off device", "device", device.DeviceID(), "error", err)
		}
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "lights turned off"})
}

func (h *LightsHandler) SetColor(w http.ResponseWriter, r *http.Request, color govee.Color, colorName string) {
	h.forEachDevice(func(device *govee.Device) {
		err := device.SetColor(color)
		if err != nil {
			h.Logger.Error("Failed to set color", "device", device.DeviceID(), "color", colorName, "error", err)
		}
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "lights set to " + colorName})
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
	var req struct {
		R int `json:"r"`
		G int `json:"g"`
		B int `json:"b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.R < 0 || req.R > 255 || req.G < 0 || req.G > 255 || req.B < 0 || req.B > 255 {
		http.Error(w, "RGB values must be between 0 and 255", http.StatusBadRequest)
		return
	}
	color := govee.Color{R: uint(req.R), G: uint(req.G), B: uint(req.B)}
	h.SetColor(w, r, color, "rgb")
}

func (h *LightsHandler) Brightness(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Brightness int `json:"brightness"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Brightness < 0 || req.Brightness > 100 {
		http.Error(w, "Brightness must be between 0 and 100", http.StatusBadRequest)
		return
	}
	h.forEachDevice(func(device *govee.Device) {
		err := device.SetBrightness(govee.Brightness(req.Brightness))
		if err != nil {
			h.Logger.Error("Failed to set brightness", "device", device.DeviceID(), "brightness", req.Brightness, "error", err)
		}
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "brightness set"})
}

func (h *LightsHandler) Status(w http.ResponseWriter, r *http.Request) {
	var statuses []map[string]interface{}
	for _, device := range h.Controller.Devices() {
		err := device.RequestStatus()
		if err != nil {
			h.Logger.Error("Failed to request status", "device", device.DeviceID(), "error", err)
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
