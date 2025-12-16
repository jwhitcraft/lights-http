package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	token := "test-token"
	middleware := AuthMiddleware(token)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer test-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer wrong-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing bearer prefix",
			authHeader:     "test-token",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "no auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
