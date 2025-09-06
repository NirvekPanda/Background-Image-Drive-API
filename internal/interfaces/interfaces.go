package interfaces

import (
	"context"
)

// DatabaseService defines the interface for database operations
type DatabaseService interface {
	// Connection management
	Close() error
	GetDB() interface{} // Returns the underlying database connection

	// Image operations
	CreateImage(ctx context.Context, image interface{}) error
	GetImage(ctx context.Context, imageID string) (interface{}, error)
	ListImages(ctx context.Context) ([]interface{}, error)
	GetImageCount(ctx context.Context) (int32, error)
	DeleteImage(ctx context.Context, imageID string) error
	GetCurrentImage(ctx context.Context) (interface{}, error)

	// Location operations
	CreateLocation(ctx context.Context, imageID string, location interface{}) error
	GetLocation(ctx context.Context, imageID string) (interface{}, error)
	UpdateLocation(ctx context.Context, imageID string, location interface{}) error
	DeleteLocation(ctx context.Context, imageID string) error
}

// ImageService defines the interface for image-related operations
type ImageService interface {
	// Core image operations
	GetCurrentImage(ctx context.Context, req interface{}) (interface{}, error)
	UploadImage(ctx context.Context, req interface{}) (interface{}, error)
	GetImageCount(ctx context.Context, req interface{}) (interface{}, error)
	ListImages(ctx context.Context, req interface{}) (interface{}, error)
	GetImageById(ctx context.Context, req interface{}) (interface{}, error)
	DeleteImage(ctx context.Context, req interface{}) (interface{}, error)
}

// LocationService defines the interface for location-related operations
type LocationService interface {
	// Location operations
	GetLocationFromCoords(ctx context.Context, req interface{}) (interface{}, error)
	GetLocationFromName(ctx context.Context, req interface{}) (interface{}, error)
}

// DriveService defines the interface for Google Drive operations
type DriveService interface {
	// File operations
	UploadFile(ctx context.Context, filename string, data []byte) (string, error)
	DeleteFile(ctx context.Context, fileID string) error
	GetFile(ctx context.Context, fileID string) ([]byte, error)
	GetFileURL(ctx context.Context, fileID string) (string, error)
}
