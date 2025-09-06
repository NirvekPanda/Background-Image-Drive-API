package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/database"
	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
)

// ImageService implements the gRPC ImageService
type ImageService struct {
	pb.UnimplementedImageServiceServer
	driveUtil *DriveUtilOAuth
	dbService *database.DatabaseService
}

// NewImageService creates a new ImageService instance
func NewImageService(driveUtil *DriveUtilOAuth, dbService *database.DatabaseService) *ImageService {
	return &ImageService{
		driveUtil: driveUtil,
		dbService: dbService,
	}
}

// GetCurrentImage returns the most recently created image
func (s *ImageService) GetCurrentImage(ctx context.Context, req *pb.GetCurrentImageRequest) (*pb.GetCurrentImageResponse, error) {
	image, err := s.dbService.GetCurrentImage(ctx)
	if err != nil {
		return &pb.GetCurrentImageResponse{
			Success: false,
			Message: "No images found",
		}, nil
	}

	return &pb.GetCurrentImageResponse{
		Success:  true,
		Message:  "Current image retrieved successfully",
		Metadata: image,
	}, nil
}

// UploadImage uploads an image to Google Drive and stores metadata
func (s *ImageService) UploadImage(ctx context.Context, req *pb.UploadImageRequest) (*pb.UploadImageResponse, error) {
	// Generate a unique ID if not provided
	imageID := req.Id
	if imageID == "" {
		imageID = fmt.Sprintf("img_%d", time.Now().UnixNano())
	}

	// Generate filename
	filename := fmt.Sprintf("%s_%s.jpg", imageID, req.Title)
	if filename == "_" {
		filename = fmt.Sprintf("%s.jpg", imageID)
	}

	// Upload to Google Drive
	driveFileID, err := s.driveUtil.UploadFile(ctx, filename, req.ImageData)
	if err != nil {
		return &pb.UploadImageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to upload to Google Drive: %v", err),
		}, nil
	}

	// Create metadata
	metadata := &pb.ImageMetadata{
		Id:          imageID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		DriveFileId: driveFileID,
	}

	// Store metadata in database
	err = s.dbService.CreateImage(ctx, metadata)
	if err != nil {
		return &pb.UploadImageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to save metadata: %v", err),
		}, nil
	}

	return &pb.UploadImageResponse{
		Success:  true,
		Message:  "Image uploaded successfully",
		Metadata: metadata,
	}, nil
}

// GetImageCount returns the total number of images
func (s *ImageService) GetImageCount(ctx context.Context, req *pb.GetImageCountRequest) (*pb.GetImageCountResponse, error) {
	count, err := s.dbService.GetImageCount(ctx)
	if err != nil {
		return &pb.GetImageCountResponse{
			Count: 0,
		}, err
	}

	return &pb.GetImageCountResponse{
		Count: count,
	}, nil
}

// ListImages returns all images from the database
func (s *ImageService) ListImages(ctx context.Context, req *pb.ListImagesRequest) (*pb.ListImagesResponse, error) {
	images, err := s.dbService.ListImages(ctx)
	if err != nil {
		return &pb.ListImagesResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to list images: %v", err),
		}, nil
	}

	return &pb.ListImagesResponse{
		Success: true,
		Message: fmt.Sprintf("Found %d images", len(images)),
		Images:  images,
	}, nil
}

// GetImageById retrieves a specific image by ID
func (s *ImageService) GetImageById(ctx context.Context, req *pb.GetImageByIdRequest) (*pb.GetImageByIdResponse, error) {
	image, err := s.dbService.GetImage(ctx, req.ImageId)
	if err != nil {
		return &pb.GetImageByIdResponse{
			Success: false,
			Message: "Image not found",
		}, nil
	}

	return &pb.GetImageByIdResponse{
		Success:  true,
		Message:  "Image retrieved successfully",
		Metadata: image,
	}, nil
}

// DeleteImage removes an image from both database and Google Drive
func (s *ImageService) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*pb.DeleteImageResponse, error) {
	// Get image metadata first to get the Drive file ID
	image, err := s.dbService.GetImage(ctx, req.ImageId)
	if err != nil {
		return &pb.DeleteImageResponse{
			Success: false,
			Message: "Image not found",
		}, nil
	}

	// Delete from Google Drive
	if image.DriveFileId != "" {
		err := s.driveUtil.DeleteFile(ctx, image.DriveFileId)
		if err != nil {
			return &pb.DeleteImageResponse{
				Success: false,
				Message: fmt.Sprintf("Failed to delete from Google Drive: %v", err),
			}, nil
		}
	}

	// Delete from database
	err = s.dbService.DeleteImage(ctx, req.ImageId)
	if err != nil {
		return &pb.DeleteImageResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete from database: %v", err),
		}, nil
	}

	return &pb.DeleteImageResponse{
		Success: true,
		Message: "Image deleted successfully",
	}, nil
}
