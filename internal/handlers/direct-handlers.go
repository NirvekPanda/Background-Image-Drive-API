package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/services"
	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
)

// DirectHTTPHandler contains the services directly (no gRPC)
type DirectHTTPHandler struct {
	imageService    *services.ImageService
	locationService *services.LocationService
}

// NewDirectHTTPHandler creates a new HTTP handler with direct service access
func NewDirectHTTPHandler(imageService *services.ImageService, locationService *services.LocationService) *DirectHTTPHandler {
	return &DirectHTTPHandler{
		imageService:    imageService,
		locationService: locationService,
	}
}

// RegisterRoutes sets up all HTTP routes
func (h *DirectHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	// Image endpoints
	mux.HandleFunc("GET /api/v1/images/current", h.getCurrentImage)
	mux.HandleFunc("POST /api/v1/images/upload", h.uploadImage)
	mux.HandleFunc("GET /api/v1/images/count", h.getImageCount)
	mux.HandleFunc("GET /api/v1/images", h.listImages)
	mux.HandleFunc("GET /api/v1/images/{id}", h.getImageById)
	mux.HandleFunc("DELETE /api/v1/images/{id}", h.deleteImage)

	// Location endpoints
	mux.HandleFunc("GET /api/v1/location/coords", h.getLocationFromCoords)
	mux.HandleFunc("GET /api/v1/location/name", h.getLocationFromName)

	// Health check
	mux.HandleFunc("GET /health", h.healthCheck)
}

// GET /api/v1/images/current
func (h *DirectHTTPHandler) getCurrentImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	image, err := h.imageService.GetCurrentImage(ctx, &pb.GetCurrentImageRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get current image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(image)
}

// POST /api/v1/images/upload
func (h *DirectHTTPHandler) uploadImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Parse multipart form (10MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	imageData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read image data", http.StatusInternalServerError)
		return
	}

	// Get form values
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Parse location data
	location := &pb.Location{}

	// Try to get coordinates
	if latStr := r.FormValue("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			location.Latitude = lat
		}
	}

	if lngStr := r.FormValue("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			location.Longitude = lng
		}
	}

	// Get other location fields
	location.Name = r.FormValue("location_name")
	location.Country = r.FormValue("country")
	location.City = r.FormValue("city")
	location.Address = r.FormValue("address")

	// Create upload request
	req := &pb.UploadImageRequest{
		Title:       title,
		Description: description,
		Location:    location,
		ImageData:   imageData,
	}

	// Call service directly
	resp, err := h.imageService.UploadImage(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /api/v1/images/count
func (h *DirectHTTPHandler) getImageCount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.imageService.GetImageCount(ctx, &pb.GetImageCountRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get image count: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /api/v1/images
func (h *DirectHTTPHandler) listImages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.imageService.ListImages(ctx, &pb.ListImagesRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list images: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /api/v1/images/{id}
func (h *DirectHTTPHandler) getImageById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	imageId := r.PathValue("id")
	if imageId == "" {
		http.Error(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	req := &pb.GetImageByIdRequest{
		ImageId: imageId,
	}

	resp, err := h.imageService.GetImageById(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DELETE /api/v1/images/{id}
func (h *DirectHTTPHandler) deleteImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	imageId := r.PathValue("id")
	if imageId == "" {
		http.Error(w, "Image ID is required", http.StatusBadRequest)
		return
	}

	req := &pb.DeleteImageRequest{
		ImageId: imageId,
	}

	resp, err := h.imageService.DeleteImage(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /api/v1/location/coords?lat=37.7749&lng=-122.4194
func (h *DirectHTTPHandler) getLocationFromCoords(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")

	if latStr == "" || lngStr == "" {
		http.Error(w, "Both lat and lng parameters are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Invalid latitude value", http.StatusBadRequest)
		return
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		http.Error(w, "Invalid longitude value", http.StatusBadRequest)
		return
	}

	req := &pb.GetLocationFromCoordsRequest{
		Latitude:  lat,
		Longitude: lng,
	}

	resp, err := h.locationService.GetLocationFromCoords(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get location: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /api/v1/location/name?name=San Francisco
func (h *DirectHTTPHandler) getLocationFromName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	locationName := r.URL.Query().Get("name")
	if locationName == "" {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	req := &pb.GetLocationFromNameRequest{
		LocationName: locationName,
	}

	resp, err := h.locationService.GetLocationFromName(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get location: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GET /health
func (h *DirectHTTPHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339Nano),
		"service":   "image-api-cloudrun",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
