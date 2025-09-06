package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteDatabase creates a SQLite database service for local development
func NewSQLiteDatabase(ctx context.Context) (*DatabaseService, error) {
	// Use in-memory database for development
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		dbPath = ":memory:" // In-memory database
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	// Test the connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create tables
	if err := createSQLiteTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite doesn't support multiple connections well
	db.SetMaxIdleConns(1)

	return &DatabaseService{db: db}, nil
}

// createSQLiteTables creates the necessary tables for SQLite
func createSQLiteTables(db *sql.DB) error {
	// Create images table
	imagesTable := `
		CREATE TABLE IF NOT EXISTS images (
			id TEXT PRIMARY KEY,
			title TEXT,
			description TEXT,
			drive_file_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`

	// Create locations table
	locationsTable := `
		CREATE TABLE IF NOT EXISTS locations (
			image_id TEXT PRIMARY KEY,
			latitude REAL,
			longitude REAL,
			name TEXT,
			country TEXT,
			city TEXT,
			address TEXT,
			FOREIGN KEY (image_id) REFERENCES images (id) ON DELETE CASCADE
		)
	`

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_images_created_at ON images(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_locations_image_id ON locations(image_id)",
	}

	// Execute table creation
	if _, err := db.Exec(imagesTable); err != nil {
		return fmt.Errorf("failed to create images table: %v", err)
	}

	if _, err := db.Exec(locationsTable); err != nil {
		return fmt.Errorf("failed to create locations table: %v", err)
	}

	// Execute index creation
	for _, index := range indexes {
		if _, err := db.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
	}

	return nil
}
