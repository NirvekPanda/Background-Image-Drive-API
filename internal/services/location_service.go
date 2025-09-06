package services

import (
	"context"
	"fmt"
	"math"
	"strings"

	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	"googlemaps.github.io/maps"
)

// LocationService implements the gRPC LocationService
type LocationService struct {
	pb.UnimplementedLocationServiceServer
	mapsClient *maps.Client
}

// NewLocationService creates a new LocationService instance
func NewLocationService(apiKey string) (*LocationService, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Maps client: %v", err)
	}

	return &LocationService{
		mapsClient: client,
	}, nil
}

// GetLocationFromCoords converts coordinates to location data
func (s *LocationService) GetLocationFromCoords(ctx context.Context, req *pb.GetLocationFromCoordsRequest) (*pb.GetLocationFromCoordsResponse, error) {
	// Validate coordinates
	if !isValidLatitude(req.Latitude) || !isValidLongitude(req.Longitude) {
		return &pb.GetLocationFromCoordsResponse{
			Success: false,
			Message: "Invalid coordinates provided",
		}, nil
	}

	// Use Google Maps API for reverse geocoding
	location, err := s.reverseGeocode(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return &pb.GetLocationFromCoordsResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get location: %v", err),
		}, nil
	}

	return &pb.GetLocationFromCoordsResponse{
		Success:  true,
		Message:  "Location retrieved successfully",
		Location: location,
	}, nil
}

// GetLocationFromName converts location name to structured location data
func (s *LocationService) GetLocationFromName(ctx context.Context, req *pb.GetLocationFromNameRequest) (*pb.GetLocationFromNameResponse, error) {
	if strings.TrimSpace(req.LocationName) == "" {
		return &pb.GetLocationFromNameResponse{
			Success: false,
			Message: "Location name cannot be empty",
		}, nil
	}

	// Use Google Maps API for geocoding
	location, err := s.geocode(ctx, req.LocationName)
	if err != nil {
		return &pb.GetLocationFromNameResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get location: %v", err),
		}, nil
	}

	return &pb.GetLocationFromNameResponse{
		Success:  true,
		Message:  "Location retrieved successfully",
		Location: location,
	}, nil
}

// reverseGeocode converts coordinates to location data using Google Maps API
func (s *LocationService) reverseGeocode(ctx context.Context, lat, lng float64) (*pb.Location, error) {
	// Create a LatLng for the request
	latlng := &maps.LatLng{
		Lat: lat,
		Lng: lng,
	}

	// Make reverse geocoding request
	req := &maps.GeocodingRequest{
		LatLng: latlng,
	}

	resp, err := s.mapsClient.ReverseGeocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("reverse geocoding failed: %v", err)
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("no results found for coordinates")
	}

	// Use the first result (most accurate)
	result := resp[0]
	location := &pb.Location{
		Latitude:  lat,
		Longitude: lng,
		Address:   result.FormattedAddress,
	}

	// Extract location components
	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			switch componentType {
			case "locality":
				location.City = component.LongName
			case "country":
				location.Country = component.LongName
			case "administrative_area_level_1":
				if location.City == "" {
					location.City = component.LongName
				}
			}
		}
	}

	// Set a default name if not found
	if location.City == "" {
		location.City = "Unknown City"
	}
	if location.Country == "" {
		location.Country = "Unknown Country"
	}
	if location.Name == "" {
		location.Name = location.City
	}

	return location, nil
}

// geocode converts location name to structured location data using Google Maps API
func (s *LocationService) geocode(ctx context.Context, locationName string) (*pb.Location, error) {
	// Make geocoding request
	req := &maps.GeocodingRequest{
		Address: locationName,
	}

	resp, err := s.mapsClient.Geocode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %v", err)
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("no results found for location: %s", locationName)
	}

	// Use the first result (most accurate)
	result := resp[0]
	location := &pb.Location{
		Name:      locationName,
		Latitude:  result.Geometry.Location.Lat,
		Longitude: result.Geometry.Location.Lng,
		Address:   result.FormattedAddress,
	}

	// Extract location components
	for _, component := range result.AddressComponents {
		for _, componentType := range component.Types {
			switch componentType {
			case "locality":
				location.City = component.LongName
			case "country":
				location.Country = component.LongName
			case "administrative_area_level_1":
				if location.City == "" {
					location.City = component.LongName
				}
			}
		}
	}

	// Set defaults if not found
	if location.City == "" {
		location.City = "Unknown City"
	}
	if location.Country == "" {
		location.Country = "Unknown Country"
	}

	return location, nil
}

// isValidLatitude checks if latitude is valid (-90 to 90)
func isValidLatitude(lat float64) bool {
	return !math.IsNaN(lat) && lat >= -90.0 && lat <= 90.0
}

// isValidLongitude checks if longitude is valid (-180 to 180)
func isValidLongitude(lng float64) bool {
	return !math.IsNaN(lng) && lng >= -180.0 && lng <= 180.0
}
