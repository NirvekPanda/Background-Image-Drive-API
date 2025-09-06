package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create CORS middleware with test config
	config := &CORSConfig{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           3600,
	}

	corsHandler := CORS(config)(handler)

	t.Run("preflight_request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/v1/images/current", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "GET")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		w := httptest.NewRecorder()
		corsHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check CORS headers
		expectedHeaders := map[string]string{
			"Access-Control-Allow-Origin":      "http://localhost:3000",
			"Access-Control-Allow-Methods":     "GET, POST, OPTIONS",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization",
			"Access-Control-Allow-Credentials": "true",
		}

		for header, expectedValue := range expectedHeaders {
			actualValue := w.Header().Get(header)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s: %s, got: %s", header, expectedValue, actualValue)
			}
		}
	})

	t.Run("actual_request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/images/current", nil)
		req.Header.Set("Origin", "http://localhost:3000")

		w := httptest.NewRecorder()
		corsHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check that CORS headers are set
		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "http://localhost:3000" {
			t.Errorf("Expected Access-Control-Allow-Origin: http://localhost:3000, got: %s", origin)
		}
	})

	t.Run("disallowed_origin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/images/current", nil)
		req.Header.Set("Origin", "http://malicious-site.com")

		w := httptest.NewRecorder()
		corsHandler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Check that no CORS origin header is set for disallowed origins
		origin := w.Header().Get("Access-Control-Allow-Origin")
		if origin != "" {
			t.Errorf("Expected no Access-Control-Allow-Origin for disallowed origin, got: %s", origin)
		}
	})
}

func TestGetCORSConfig(t *testing.T) {
	// Test default configuration
	config := GetCORSConfig()

	if len(config.AllowedOrigins) == 0 {
		t.Error("Expected at least one allowed origin")
	}

	if len(config.AllowedMethods) == 0 {
		t.Error("Expected at least one allowed method")
	}

	if !config.AllowCredentials {
		t.Error("Expected AllowCredentials to be true")
	}
}
