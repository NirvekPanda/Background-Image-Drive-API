package services

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// DriveUtil handles Google Drive operations
type DriveUtil struct {
	service *drive.Service
}

// NewDriveUtil creates a new Google Drive utility client
func NewDriveUtil(ctx context.Context, credentialsJSON []byte) (*DriveUtil, error) {
	config, err := google.JWTConfigFromJSON(credentialsJSON, drive.DriveFileScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %v", err)
	}

	client := config.Client(ctx)
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Drive service: %v", err)
	}

	return &DriveUtil{service: service}, nil
}

// UploadFile uploads image data to Google Drive and returns the file ID
func (d *DriveUtil) UploadFile(ctx context.Context, filename string, imageData []byte) (string, error) {
	file := &drive.File{Name: filename}

	call := d.service.Files.Create(file).Media(bytes.NewReader(imageData)).Context(ctx)
	createdFile, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	return createdFile.Id, nil
}

// DownloadFile downloads file content from Google Drive by file ID
func (d *DriveUtil) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	call := d.service.Files.Get(fileID).Context(ctx)
	resp, err := call.Download()
	if err != nil {
		return nil, fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %v", err)
	}

	return data, nil
}

// DeleteFile removes a file from Google Drive by file ID
func (d *DriveUtil) DeleteFile(ctx context.Context, fileID string) error {
	err := d.service.Files.Delete(fileID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}
