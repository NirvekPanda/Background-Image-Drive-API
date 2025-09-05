package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	"google.golang.org/grpc"
)

// HTTPHandler contains the gRPC client connections
type HTTPHandler struct {
	imageClient    pb.ImageServiceClient
	locationClient pb.LocationServiceClient
}

// NewHTTPHandler creates a new HTTP handler with gRPC clients
func NewHTTPHandler(imageConn, locationConn *grpc.ClientConn) *HTTPHandler {
	return &HTTPHandler{
		imageClient:    pb.NewImageServiceClient(imageConn),
		locationClient: pb.NewLocationServiceClient(locationConn),
	}
}

// RegisterRoutes sets up all HTTP routes
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	// Image endpoints
	mux.HandleFunc("GET /api/v1/images/current", h.getCurrentImage)
	mux.HandleFunc("POST /api/v1/images/upload", h.uploadImage)
	mux.HandleFunc("GET /api/v1/images/count", h.getImageCount)
	mux.HandleFunc("GET /api/v1/images/{id}", h.getImageById)
	mux.HandleFunc("DELETE /api/v1/images/{id}", h.deleteImage)

	// Location endpoints
	mux.HandleFunc("GET /api/v1/location/coords", h.getLocationFromCoords)
	mux.HandleFunc("GET /api/v1/location/name", h.getLocationFromName)

	// Health check
	mux.HandleFunc("GET /health", h.healthCheck)
}

// GET /api/v1/images/current
func (h *HTTPHandler) getCurrentImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.imageClient.GetCurrentImage(ctx, &pb.GetCurrentImageRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get current image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// POST /api/v1/images/upload
func (h *HTTPHandler) uploadImage(w http.ResponseWriter, r *http.Request) {
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

	// Call gRPC service
	resp, err := h.imageClient.UploadImage(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/images/count
func (h *HTTPHandler) getImageCount(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.imageClient.GetImageCount(ctx, &pb.GetImageCountRequest{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get image count: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/images/{id}
func (h *HTTPHandler) getImageById(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.imageClient.GetImageById(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DELETE /api/v1/images/{id}
func (h *HTTPHandler) deleteImage(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.imageClient.DeleteImage(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete image: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/location/coords?lat=37.7749&lng=-122.4194
func (h *HTTPHandler) getLocationFromCoords(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.locationClient.GetLocationFromCoords(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get location: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /api/v1/location/name?name=San Francisco
func (h *HTTPHandler) getLocationFromName(w http.ResponseWriter, r *http.Request) {
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

	resp, err := h.locationClient.GetLocationFromName(ctx, req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get location: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /health
func (h *HTTPHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339Nano),
		"service":   "image-api-http-gateway",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
