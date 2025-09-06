package database

import (
	"context"
	"fmt"
	"os"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/interfaces"
)

// DatabaseType represents the type of database to use
type DatabaseType string

const (
	DatabaseTypeSQLite     DatabaseType = "sqlite"
	DatabaseTypePostgreSQL DatabaseType = "postgres"
	DatabaseTypeCloudSQL   DatabaseType = "cloudsql"
)

// NewDatabaseServiceFromEnv creates a database service based on environment or configuration
func NewDatabaseServiceFromEnv(ctx context.Context) (interfaces.DatabaseService, error) {
	// Check environment variable first
	dbType := os.Getenv("DATABASE_TYPE")
	if dbType == "" {
		// Default to SQLite for local development
		dbType = string(DatabaseTypeSQLite)
	}

	switch DatabaseType(dbType) {
	case DatabaseTypeSQLite:
		return NewSQLiteDatabase(ctx)
	case DatabaseTypePostgreSQL:
		return NewLocalPostgres(ctx)
	case DatabaseTypeCloudSQL:
		return NewCloudSQLFromEnv(ctx)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// NewDatabaseServiceWithType creates a database service with a specific type
func NewDatabaseServiceWithType(ctx context.Context, dbType DatabaseType) (interfaces.DatabaseService, error) {
	switch dbType {
	case DatabaseTypeSQLite:
		return NewSQLiteDatabase(ctx)
	case DatabaseTypePostgreSQL:
		return NewLocalPostgres(ctx)
	case DatabaseTypeCloudSQL:
		return NewCloudSQLFromEnv(ctx)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// NewDatabaseServiceWithFallback creates a database service with fallback options
func NewDatabaseServiceWithFallback(ctx context.Context, primary DatabaseType, fallback DatabaseType) (interfaces.DatabaseService, error) {
	// Try primary database first
	db, err := NewDatabaseServiceWithType(ctx, primary)
	if err == nil {
		return db, nil
	}

	// If primary fails, try fallback
	fmt.Printf("Primary database (%s) failed, trying fallback (%s): %v\n", primary, fallback, err)
	return NewDatabaseServiceWithType(ctx, fallback)
}
