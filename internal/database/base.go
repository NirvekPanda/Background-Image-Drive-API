package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/interfaces"
	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
)

// BaseDatabaseService implements the DatabaseService interface
type BaseDatabaseService struct {
	db *sql.DB
}

// NewBaseDatabaseService creates a new base database service
func NewBaseDatabaseService(db *sql.DB) interfaces.DatabaseService {
	return &BaseDatabaseService{db: db}
}

// Close closes the database connection
func (d *BaseDatabaseService) Close() error {
	return d.db.Close()
}

// GetDB returns the underlying database connection
func (d *BaseDatabaseService) GetDB() interface{} {
	return d.db
}

// CreateImage creates a new image record in the database
func (d *BaseDatabaseService) CreateImage(ctx context.Context, image interface{}) error {
	img, ok := image.(*pb.ImageMetadata)
	if !ok {
		return fmt.Errorf("invalid image type")
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		_ = tx.Rollback() // Ignore rollback error in defer
	}()

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
	_, err = tx.ExecContext(ctx, query, img.Id, img.Title, img.Description, img.DriveFileId)
	if err != nil {
		return fmt.Errorf("failed to insert image: %v", err)
	}

	// Insert location if provided
	if img.Location != nil {
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
			img.Id,
			img.Location.Latitude,
			img.Location.Longitude,
			img.Location.Name,
			img.Location.Country,
			img.Location.City,
			img.Location.Address,
		)
		if err != nil {
			return fmt.Errorf("failed to insert location: %v", err)
		}
	}

	return tx.Commit()
}

// GetImage retrieves an image by ID
func (d *BaseDatabaseService) GetImage(ctx context.Context, imageID string) (interface{}, error) {
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
func (d *BaseDatabaseService) ListImages(ctx context.Context) ([]interface{}, error) {
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

	var images []interface{}
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
func (d *BaseDatabaseService) GetImageCount(ctx context.Context) (int32, error) {
	query := "SELECT COUNT(*) FROM images"
	var count int32

	err := d.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get image count: %v", err)
	}

	return count, nil
}

// DeleteImage deletes an image and its location data
func (d *BaseDatabaseService) DeleteImage(ctx context.Context, imageID string) error {
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
func (d *BaseDatabaseService) GetCurrentImage(ctx context.Context) (interface{}, error) {
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

// CreateLocation creates a location record
func (d *BaseDatabaseService) CreateLocation(ctx context.Context, imageID string, location interface{}) error {
	loc, ok := location.(*pb.Location)
	if !ok {
		return fmt.Errorf("invalid location type")
	}

	query := `
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
	_, err := d.db.ExecContext(ctx, query,
		imageID,
		loc.Latitude,
		loc.Longitude,
		loc.Name,
		loc.Country,
		loc.City,
		loc.Address,
	)
	return err
}

// GetLocation retrieves a location by image ID
func (d *BaseDatabaseService) GetLocation(ctx context.Context, imageID string) (interface{}, error) {
	query := `
		SELECT latitude, longitude, name, country, city, address
		FROM locations
		WHERE image_id = $1
	`

	var location pb.Location
	err := d.db.QueryRowContext(ctx, query, imageID).Scan(
		&location.Latitude,
		&location.Longitude,
		&location.Name,
		&location.Country,
		&location.City,
		&location.Address,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("location not found")
		}
		return nil, fmt.Errorf("failed to get location: %v", err)
	}

	return &location, nil
}

// UpdateLocation updates a location record
func (d *BaseDatabaseService) UpdateLocation(ctx context.Context, imageID string, location interface{}) error {
	loc, ok := location.(*pb.Location)
	if !ok {
		return fmt.Errorf("invalid location type")
	}

	query := `
		UPDATE locations SET
			latitude = $2,
			longitude = $3,
			name = $4,
			country = $5,
			city = $6,
			address = $7
		WHERE image_id = $1
	`
	result, err := d.db.ExecContext(ctx, query,
		imageID,
		loc.Latitude,
		loc.Longitude,
		loc.Name,
		loc.Country,
		loc.City,
		loc.Address,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("location not found")
	}

	return nil
}

// DeleteLocation deletes a location record
func (d *BaseDatabaseService) DeleteLocation(ctx context.Context, imageID string) error {
	query := "DELETE FROM locations WHERE image_id = $1"
	_, err := d.db.ExecContext(ctx, query, imageID)
	return err
}
