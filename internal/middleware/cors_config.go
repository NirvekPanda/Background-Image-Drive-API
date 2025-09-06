package middleware

import (
	"os"
	"strings"
)

// GetCORSConfig returns CORS configuration based on environment variables
func GetCORSConfig() *CORSConfig {
	config := &CORSConfig{
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"HEAD",
		},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Language",
			"Content-Language",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"Origin",
		},
		ExposedHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           86400, // 24 hours
	}

	// Get allowed origins from environment variable
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOriginsEnv != "" {
		// Split by comma and trim spaces
		origins := strings.Split(allowedOriginsEnv, ",")
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
		config.AllowedOrigins = origins
	} else {
		// Default origins for development
		config.AllowedOrigins = []string{
			"http://localhost:3000",    // React dev server
			"http://localhost:3001",    // Alternative React port
			"http://localhost:3002",    // Additional dev port
			"http://localhost:3003",    // Additional dev port
			"http://127.0.0.1:3000",    // Alternative localhost
			"http://127.0.0.1:3001",    // Alternative localhost
			"http://127.0.0.1:3002",    // Alternative localhost
			"http://127.0.0.1:3003",    // Alternative localhost
			"https://nirvekpandey.com", // Production domain
		}
	}

	// Allow all origins in development (be careful in production!)
	if os.Getenv("ENVIRONMENT") == "development" || os.Getenv("ENV") == "dev" {
		config.AllowedOrigins = []string{"*"}
	}

	return config
}
