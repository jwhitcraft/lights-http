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

package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	Logger *slog.Logger
}

func (m *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := generateRequestID()
		ctx := context.WithValue(r.Context(), "requestID", requestID)
		r = r.WithContext(ctx)

		// Set request ID in response header
		w.Header().Set("X-Request-ID", requestID)

		// Log request
		m.Logger.Info("Request started",
			"method", r.Method,
			"path", r.URL.Path,
			"requestID", requestID,
			"userAgent", r.UserAgent(),
			"remoteAddr", r.RemoteAddr,
		)

		// Wrap response writer to capture status
		wrapped := &responseWriter{ResponseWriter: w, status: 200}

		next.ServeHTTP(wrapped, r)

		// Log response
		duration := time.Since(start)
		m.Logger.Info("Request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"requestID", requestID,
			"status", wrapped.status,
			"duration", duration,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// generateRequestID creates a simple 8-character hex request ID
func generateRequestID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}