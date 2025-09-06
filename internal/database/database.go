package database

import (
	"context"
	"fmt"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/interfaces"
)

// LegacyDatabaseService is a compatibility wrapper for the old DatabaseService
// This maintains backward compatibility while transitioning to the new interface-based system
type LegacyDatabaseService struct {
	service interfaces.DatabaseService
}

// NewLegacyDatabaseService creates a legacy database service wrapper
func NewLegacyDatabaseService(ctx context.Context) (*LegacyDatabaseService, error) {
	service, err := NewDatabaseServiceFromEnv(ctx)
	if err != nil {
		return nil, err
	}
	return &LegacyDatabaseService{service: service}, nil
}

// Close closes the database connection
func (d *LegacyDatabaseService) Close() error {
	return d.service.Close()
}

// GetDB returns the underlying database connection
func (d *LegacyDatabaseService) GetDB() interface{} {
	return d.service.GetDB()
}

// CreateImage creates a new image record in the database
func (d *LegacyDatabaseService) CreateImage(ctx context.Context, image interface{}) error {
	return d.service.CreateImage(ctx, image)
}

// GetImage retrieves an image by ID
func (d *LegacyDatabaseService) GetImage(ctx context.Context, imageID string) (interface{}, error) {
	return d.service.GetImage(ctx, imageID)
}

// ListImages retrieves all images
func (d *LegacyDatabaseService) ListImages(ctx context.Context) ([]interface{}, error) {
	return d.service.ListImages(ctx)
}

// GetImageCount returns the total number of images
func (d *LegacyDatabaseService) GetImageCount(ctx context.Context) (int32, error) {
	return d.service.GetImageCount(ctx)
}

// DeleteImage deletes an image and its location data
func (d *LegacyDatabaseService) DeleteImage(ctx context.Context, imageID string) error {
	return d.service.DeleteImage(ctx, imageID)
}

// GetCurrentImage returns the most recently created image
func (d *LegacyDatabaseService) GetCurrentImage(ctx context.Context) (interface{}, error) {
	return d.service.GetCurrentImage(ctx)
}

// CreateLocation creates a location record
func (d *LegacyDatabaseService) CreateLocation(ctx context.Context, imageID string, location interface{}) error {
	return d.service.CreateLocation(ctx, imageID, location)
}

// GetLocation retrieves a location by image ID
func (d *LegacyDatabaseService) GetLocation(ctx context.Context, imageID string) (interface{}, error) {
	return d.service.GetLocation(ctx, imageID)
}

// UpdateLocation updates a location record
func (d *LegacyDatabaseService) UpdateLocation(ctx context.Context, imageID string, location interface{}) error {
	return d.service.UpdateLocation(ctx, imageID, location)
}

// DeleteLocation deletes a location record
func (d *LegacyDatabaseService) DeleteLocation(ctx context.Context, imageID string) error {
	return d.service.DeleteLocation(ctx, imageID)
}

// NewDatabaseServiceLegacy creates a new database service (legacy function for backward compatibility)
func NewDatabaseServiceLegacy(connectionString string) (*LegacyDatabaseService, error) {
	return nil, fmt.Errorf("use NewLegacyDatabaseService or NewDatabaseServiceWithType instead")
}
