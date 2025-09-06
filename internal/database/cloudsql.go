package database

import (
	"context"
	"fmt"
	"os"
)

// CloudSQLConfig holds Cloud SQL connection configuration
type CloudSQLConfig struct {
	InstanceConnectionName string
	DatabaseName           string
	User                   string
	Password               string
}

// NewCloudSQLConnection creates a Cloud SQL connection using direct connection
func NewCloudSQLConnection(ctx context.Context, config CloudSQLConfig) (*DatabaseService, error) {
	// Use public IP for direct connection to Cloud SQL
	// In production, you should use Cloud SQL Auth Proxy or private IP
	connectionString := fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=require",
		config.InstanceConnectionName, // This will be the IP address
		config.User,
		config.Password,
		config.DatabaseName,
	)

	return NewDatabaseService(connectionString)
}

// NewCloudSQLFromEnv creates a Cloud SQL connection from environment variables
func NewCloudSQLFromEnv(ctx context.Context) (*DatabaseService, error) {
	config := CloudSQLConfig{
		InstanceConnectionName: os.Getenv("CLOUD_SQL_CONNECTION_NAME"),
		DatabaseName:           os.Getenv("CLOUD_SQL_DATABASE"),
		User:                   os.Getenv("CLOUD_SQL_USER"),
		Password:               os.Getenv("CLOUD_SQL_PASSWORD"),
	}

	// Validate required environment variables
	if config.InstanceConnectionName == "" {
		return nil, fmt.Errorf("CLOUD_SQL_CONNECTION_NAME environment variable is required")
	}
	if config.DatabaseName == "" {
		return nil, fmt.Errorf("CLOUD_SQL_DATABASE environment variable is required")
	}
	if config.User == "" {
		return nil, fmt.Errorf("CLOUD_SQL_USER environment variable is required")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("CLOUD_SQL_PASSWORD environment variable is required")
	}

	return NewCloudSQLConnection(ctx, config)
}

// NewLocalPostgres creates a local PostgreSQL connection for development
func NewLocalPostgres(ctx context.Context) (*DatabaseService, error) {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		// Default local connection string
		connectionString = "host=localhost port=5432 user=postgres password=postgres dbname=portfolio_images sslmode=disable"
	}

	return NewDatabaseService(connectionString)
}
