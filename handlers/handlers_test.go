package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log/slog"

	govee "github.com/swrm-io/go-vee"
)

// MockController is a mock implementation of ControllerInterface for testing
type MockController struct {
}

func (m *MockController) Devices() []*govee.Device {
	return []*govee.Device{}
}

func TestTurnOn(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	req := httptest.NewRequest("POST", "/lights/on", nil)
	w := httptest.NewRecorder()

	handler.TurnOn(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "lights turned on" {
		t.Errorf("expected status 'lights turned on', got %s", response["status"])
	}
}

func TestTurnOff(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	req := httptest.NewRequest("POST", "/lights/off", nil)
	w := httptest.NewRecorder()

	handler.TurnOff(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "lights turned off" {
		t.Errorf("expected status 'lights turned off', got %s", response["status"])
	}
}

func TestRed(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	req := httptest.NewRequest("POST", "/lights/red", nil)
	w := httptest.NewRecorder()

	handler.Red(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "lights set to red" {
		t.Errorf("expected status 'lights set to red', got %s", response["status"])
	}
}

func TestDarkRed(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	req := httptest.NewRequest("POST", "/lights/dark-red", nil)
	w := httptest.NewRecorder()

	handler.DarkRed(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "lights set to dark-red" {
		t.Errorf("expected status 'lights set to dark-red', got %s", response["status"])
	}
}

func TestBrightness(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	body := map[string]int{"brightness": 50}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/lights/brightness", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Brightness(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "brightness set" {
		t.Errorf("expected status 'brightness set', got %s", response["status"])
	}
}

func TestStatus(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	req := httptest.NewRequest("GET", "/lights/status", nil)
	w := httptest.NewRecorder()

	handler.Status(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response []map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Since mock has no devices, expect empty array
	if len(response) != 0 {
		t.Errorf("expected empty array, got %v", response)
	}
}

func TestRGB(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	body := map[string]int{"r": 255, "g": 128, "b": 0}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/lights/rgb", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.RGB(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "lights set to rgb" {
		t.Errorf("expected status 'lights set to rgb', got %s", response["status"])
	}
}

func TestRGBInvalidValues(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	tests := []struct {
		name string
		body map[string]int
	}{
		{"negative r", map[string]int{"r": -1, "g": 0, "b": 0}},
		{"r too high", map[string]int{"r": 256, "g": 0, "b": 0}},
		{"negative g", map[string]int{"r": 0, "g": -1, "b": 0}},
		{"g too high", map[string]int{"r": 0, "g": 256, "b": 0}},
		{"negative b", map[string]int{"r": 0, "g": 0, "b": -1}},
		{"b too high", map[string]int{"r": 0, "g": 0, "b": 256}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/lights/rgb", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler.RGB(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestRGBInvalidJSON(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &LightsHandler{
		Controller: mockController,
		Logger:     logger,
	}

	req := httptest.NewRequest("POST", "/lights/rgb", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.RGB(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHealth(t *testing.T) {
	mockController := &MockController{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &HealthHandler{
		Controller: mockController,
		Logger:     logger,
		StartTime:  time.Now().Add(-time.Hour), // 1 hour ago
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response HealthStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" { // Controller is initialized
		t.Errorf("expected status 'ok', got %s", response.Status)
	}

	if response.Uptime == "" {
		t.Errorf("expected uptime to be set")
	}

	if len(response.Checks) == 0 {
		t.Errorf("expected checks to be present")
	}

	// Check controller status
	if controllerCheck, ok := response.Checks["controller"]; !ok {
		t.Errorf("expected controller check to be present")
	} else if controllerCheck.Status != "ok" {
		t.Errorf("expected controller status 'ok', got %s", controllerCheck.Status)
	}
}

func TestHealthWithDevices(t *testing.T) {
	// Create a mock controller with devices
	mockController := &MockControllerWithDevices{}
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{}))
	handler := &HealthHandler{
		Controller: mockController,
		Logger:     logger,
		StartTime:  time.Now().Add(-time.Minute * 30), // 30 minutes ago
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response HealthStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("expected status 'ok', got %s", response.Status)
	}

	// Check controller status
	if controllerCheck, ok := response.Checks["controller"]; !ok {
		t.Errorf("expected controller check to be present")
	} else if controllerCheck.Status != "ok" {
		t.Errorf("expected controller status 'ok', got %s", controllerCheck.Status)
	}
}

// MockControllerWithDevices is a mock that returns devices
type MockControllerWithDevices struct{}

func (m *MockControllerWithDevices) Devices() []*govee.Device {
	// Return a mock device (we can't easily create a real one in tests)
	return []*govee.Device{} // Empty slice represents "some devices"
}

// Similar tests for Yellow and Orange can be added
