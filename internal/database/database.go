package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	_ "github.com/lib/pq"
)

// DatabaseService handles database operations
type DatabaseService struct {
	db *sql.DB
}

// NewDatabaseService creates a new database service
func NewDatabaseService(connectionString string) (*DatabaseService, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &DatabaseService{db: db}, nil
}

// Close closes the database connection
func (d *DatabaseService) Close() error {
	return d.db.Close()
}

// GetDB returns the underlying database connection
func (d *DatabaseService) GetDB() *sql.DB {
	return d.db
}

// CreateImage creates a new image record in the database
func (d *DatabaseService) CreateImage(ctx context.Context, image *pb.ImageMetadata) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Insert image
	query := `
		INSERT INTO images (id, title, description, drive_file_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			drive_file_id = EXCLUDED.drive_file_id,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err = tx.ExecContext(ctx, query, image.Id, image.Title, image.Description, image.DriveFileId)
	if err != nil {
		return fmt.Errorf("failed to insert image: %v", err)
	}

	// Insert location if provided
	if image.Location != nil {
		locationQuery := `
			INSERT INTO locations (image_id, latitude, longitude, name, country, city, address)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (image_id) DO UPDATE SET
				latitude = EXCLUDED.latitude,
				longitude = EXCLUDED.longitude,
				name = EXCLUDED.name,
				country = EXCLUDED.country,
				city = EXCLUDED.city,
				address = EXCLUDED.address
		`
		_, err = tx.ExecContext(ctx, locationQuery,
			image.Id,
			image.Location.Latitude,
			image.Location.Longitude,
			image.Location.Name,
			image.Location.Country,
			image.Location.City,
			image.Location.Address,
		)
		if err != nil {
			return fmt.Errorf("failed to insert location: %v", err)
		}
	}

	return tx.Commit()
}

// GetImage retrieves an image by ID
func (d *DatabaseService) GetImage(ctx context.Context, imageID string) (*pb.ImageMetadata, error) {
	query := `
		SELECT i.id, i.title, i.description, i.drive_file_id, i.created_at,
		       l.latitude, l.longitude, l.name, l.country, l.city, l.address
		FROM images i
		LEFT JOIN locations l ON i.id = l.image_id
		WHERE i.id = $1
	`

	var image pb.ImageMetadata
	var location pb.Location
	var createdAt time.Time

	err := d.db.QueryRowContext(ctx, query, imageID).Scan(
		&image.Id,
		&image.Title,
		&image.Description,
		&image.DriveFileId,
		&createdAt,
		&location.Latitude,
		&location.Longitude,
		&location.Name,
		&location.Country,
		&location.City,
		&location.Address,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to get image: %v", err)
	}

	// Only set location if it has data
	if location.Latitude != 0 || location.Longitude != 0 || location.Name != "" {
		image.Location = &location
	}

	return &image, nil
}

// ListImages retrieves all images
func (d *DatabaseService) ListImages(ctx context.Context) ([]*pb.ImageMetadata, error) {
	query := `
		SELECT i.id, i.title, i.description, i.drive_file_id, i.created_at,
		       l.latitude, l.longitude, l.name, l.country, l.city, l.address
		FROM images i
		LEFT JOIN locations l ON i.id = l.image_id
		ORDER BY i.created_at DESC
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %v", err)
	}
	defer rows.Close()

	var images []*pb.ImageMetadata
	for rows.Next() {
		var image pb.ImageMetadata
		var location pb.Location
		var createdAt time.Time

		err := rows.Scan(
			&image.Id,
			&image.Title,
			&image.Description,
			&image.DriveFileId,
			&createdAt,
			&location.Latitude,
			&location.Longitude,
			&location.Name,
			&location.Country,
			&location.City,
			&location.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan image: %v", err)
		}

		// Only set location if it has data
		if location.Latitude != 0 || location.Longitude != 0 || location.Name != "" {
			image.Location = &location
		}

		images = append(images, &image)
	}

	return images, nil
}

// GetImageCount returns the total number of images
func (d *DatabaseService) GetImageCount(ctx context.Context) (int32, error) {
	query := "SELECT COUNT(*) FROM images"
	var count int32

	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get image count: %v", err)
	}

	return count, nil
}

// DeleteImage deletes an image and its location data
func (d *DatabaseService) DeleteImage(ctx context.Context, imageID string) error {
	query := "DELETE FROM images WHERE id = $1"
	result, err := d.db.ExecContext(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

// GetCurrentImage returns the most recently created image
func (d *DatabaseService) GetCurrentImage(ctx context.Context) (*pb.ImageMetadata, error) {
	query := `
		SELECT i.id, i.title, i.description, i.drive_file_id, i.created_at,
		       l.latitude, l.longitude, l.name, l.country, l.city, l.address
		FROM images i
		LEFT JOIN locations l ON i.id = l.image_id
		ORDER BY i.created_at DESC
		LIMIT 1
	`

	var image pb.ImageMetadata
	var location pb.Location
	var createdAt time.Time

	err := d.db.QueryRowContext(ctx, query).Scan(
		&image.Id,
		&image.Title,
		&image.Description,
		&image.DriveFileId,
		&createdAt,
		&location.Latitude,
		&location.Longitude,
		&location.Name,
		&location.Country,
		&location.City,
		&location.Address,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no images found")
		}
		return nil, fmt.Errorf("failed to get current image: %v", err)
	}

	// Only set location if it has data
	if location.Latitude != 0 || location.Longitude != 0 || location.Name != "" {
		image.Location = &location
	}

	return &image, nil
}
