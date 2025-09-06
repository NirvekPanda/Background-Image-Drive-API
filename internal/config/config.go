package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	// Server configuration
	GRPCPort string
	HTTPPort string

	// Google Drive configuration
	GoogleDriveCredentialsPath string
	GoogleDriveFolderID        string

	// Google Maps configuration
	GoogleMapsAPIKey string

	// Cloud SQL configuration
	CloudSQLConnectionName string
	CloudSQLDatabase       string
	CloudSQLUser           string
	CloudSQLPassword       string

	// Database configuration
	DatabaseURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		GRPCPort:                   getEnv("GRPC_PORT", "50051"),
		HTTPPort:                   getEnv("HTTP_PORT", "8080"),
		GoogleDriveCredentialsPath: getEnv("GOOGLE_DRIVE_CREDENTIALS_PATH", "credentials.json"),
		GoogleDriveFolderID:        getEnv("GOOGLE_DRIVE_FOLDER_ID", ""),
		GoogleMapsAPIKey:           getEnv("GOOGLE_MAPS_API_KEY", ""),
		CloudSQLConnectionName:     getEnv("CLOUD_SQL_CONNECTION_NAME", ""),
		CloudSQLDatabase:           getEnv("CLOUD_SQL_DATABASE", ""),
		CloudSQLUser:               getEnv("CLOUD_SQL_USER", ""),
		CloudSQLPassword:           getEnv("CLOUD_SQL_PASSWORD", ""),
		DatabaseURL:                getEnv("DATABASE_URL", ""),
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
