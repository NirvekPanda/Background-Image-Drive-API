package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/database"
	"github.com/NirvekPanda/Background-Image-Drive-API/internal/services"
	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get Google Drive folder ID
	folderID := os.Getenv("GOOGLE_DRIVE_FOLDER_ID")
	if folderID == "" {
		log.Fatalf("GOOGLE_DRIVE_FOLDER_ID environment variable is required")
	}

	// Create Google Drive utility with OAuth2
	ctx := context.Background()

	// Check for OAuth2 credentials in secrets directory (production) or current directory (local)
	oauthConfigPath := "/app/secrets/oauth_credentials.json"
	if _, err := os.Stat(oauthConfigPath); os.IsNotExist(err) {
		oauthConfigPath = "oauth_credentials.json" // Fallback to local development
	}

	tokenPath := "/app/secrets/token.json"
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		tokenPath = "token.json" // Fallback to local development
	}

	driveUtil, err := services.NewDriveUtilOAuth(ctx, oauthConfigPath, tokenPath, folderID)
	if err != nil {
		log.Fatalf("Failed to create Drive utility: %v", err)
	}

	// Get Google Maps API key
	mapsAPIKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if mapsAPIKey == "" {
		log.Fatalf("GOOGLE_MAPS_API_KEY environment variable is required")
	}

	// Create database service
	dbService, err := database.NewCloudSQLFromEnv(ctx)
	if err != nil {
		log.Printf("Failed to connect to Cloud SQL, trying local database: %v", err)
		// Fallback to local database for development
		dbService, err = database.NewLocalPostgres(ctx)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
	defer dbService.Close()

	// Create services
	imageService := services.NewImageService(driveUtil, dbService)
	locationService, err := services.NewLocationService(mapsAPIKey)
	if err != nil {
		log.Fatalf("Failed to create location service: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterImageServiceServer(grpcServer, imageService)
	pb.RegisterLocationServiceServer(grpcServer, locationService)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Start the server
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	fmt.Printf("gRPC server starting on :%s\n", port)
	fmt.Println("Services registered:")
	fmt.Println("  - ImageService")
	fmt.Println("  - LocationService")
	fmt.Println("  - gRPC Reflection enabled")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
