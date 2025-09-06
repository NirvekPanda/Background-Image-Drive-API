package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/config"
	"github.com/NirvekPanda/Background-Image-Drive-API/internal/database"
	"github.com/NirvekPanda/Background-Image-Drive-API/internal/handlers"
	"github.com/NirvekPanda/Background-Image-Drive-API/internal/middleware"
	"github.com/NirvekPanda/Background-Image-Drive-API/internal/services"
	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	ctx := context.Background()

	// Initialize Secret Manager
	secretManager, err := config.NewSecretManager(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize Secret Manager, using environment variables: %v", err)
		secretManager = nil
	}

	// Get configuration from Secret Manager or environment variables
	var folderID, mapsAPIKey, cloudSQLConnName, cloudSQLDatabase, cloudSQLUser, cloudSQLPassword, grpcPort string

	if secretManager != nil {
		// Try to get from Secret Manager first
		folderID = secretManager.GetSecretWithFallback("GOOGLE_DRIVE_FOLDER_ID", "GOOGLE_DRIVE_FOLDER_ID")
		mapsAPIKey = secretManager.GetSecretWithFallback("GOOGLE_MAPS_API_KEY", "GOOGLE_MAPS_API_KEY")
		cloudSQLConnName = secretManager.GetSecretWithFallback("CLOUD_SQL_CONNECTION_NAME", "CLOUD_SQL_CONNECTION_NAME")
		cloudSQLDatabase = secretManager.GetSecretWithFallback("CLOUD_SQL_DATABASE", "CLOUD_SQL_DATABASE")
		cloudSQLUser = secretManager.GetSecretWithFallback("CLOUD_SQL_USER", "CLOUD_SQL_USER")
		cloudSQLPassword = secretManager.GetSecretWithFallback("CLOUD_SQL_PASSWORD", "CLOUD_SQL_PASSWORD")
		grpcPort = secretManager.GetSecretWithFallback("GRPC_PORT", "GRPC_PORT")
	} else {
		// Fall back to environment variables
		folderID = os.Getenv("GOOGLE_DRIVE_FOLDER_ID")
		mapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
		cloudSQLConnName = os.Getenv("CLOUD_SQL_CONNECTION_NAME")
		cloudSQLDatabase = os.Getenv("CLOUD_SQL_DATABASE")
		cloudSQLUser = os.Getenv("CLOUD_SQL_USER")
		cloudSQLPassword = os.Getenv("CLOUD_SQL_PASSWORD")
		grpcPort = os.Getenv("GRPC_PORT")
	}

	// Validate required configuration
	if folderID == "" {
		log.Fatalf("GOOGLE_DRIVE_FOLDER_ID is required")
	}
	if mapsAPIKey == "" {
		log.Fatalf("GOOGLE_MAPS_API_KEY is required")
	}
	if cloudSQLConnName == "" {
		log.Fatalf("CLOUD_SQL_CONNECTION_NAME is required")
	}
	if cloudSQLDatabase == "" {
		log.Fatalf("CLOUD_SQL_DATABASE is required")
	}
	if cloudSQLUser == "" {
		log.Fatalf("CLOUD_SQL_USER is required")
	}
	if cloudSQLPassword == "" {
		log.Fatalf("CLOUD_SQL_PASSWORD is required")
	}
	if grpcPort == "" {
		grpcPort = "50051"
	}

	// Handle OAuth2 credentials
	var oauthConfigPath, tokenPath string
	
	if secretManager != nil {
		// In production, get credentials from Secret Manager
		oauthConfigData, err := secretManager.GetSecret("oauth-credentials")
		if err != nil {
			log.Fatalf("Failed to get OAuth credentials from Secret Manager: %v", err)
		}
		
		tokenData, err := secretManager.GetSecret("oauth-token")
		if err != nil {
			log.Fatalf("Failed to get OAuth token from Secret Manager: %v", err)
		}
		
		// Write credentials to files
		oauthConfigPath = "/tmp/oauth_credentials.json"
		tokenPath = "/tmp/token.json"
		
		if err := os.WriteFile(oauthConfigPath, []byte(oauthConfigData), 0600); err != nil {
			log.Fatalf("Failed to write OAuth credentials file: %v", err)
		}
		
		if err := os.WriteFile(tokenPath, []byte(tokenData), 0600); err != nil {
			log.Fatalf("Failed to write OAuth token file: %v", err)
		}
	} else {
		// In local development, use existing files
		oauthConfigPath = "/app/secrets/oauth_credentials.json"
		if _, err := os.Stat(oauthConfigPath); os.IsNotExist(err) {
			oauthConfigPath = "oauth_credentials.json" // Fallback to local development
		}

		tokenPath = "/app/secrets/token.json"
		if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
			tokenPath = "token.json" // Fallback to local development
		}
	}

	driveUtil, err := services.NewDriveUtilOAuth(ctx, oauthConfigPath, tokenPath, folderID)
	if err != nil {
		log.Fatalf("Failed to create Drive utility: %v", err)
	}

	// Create database service
	dbService, err := database.NewCloudSQLConnection(ctx, database.CloudSQLConfig{
		InstanceConnectionName: cloudSQLConnName,
		DatabaseName:           cloudSQLDatabase,
		User:                   cloudSQLUser,
		Password:               cloudSQLPassword,
	})
	if err != nil {
		log.Printf("Failed to connect to Cloud SQL, trying SQLite: %v", err)
		// Fallback to SQLite for development
		dbService, err = database.NewSQLiteDatabase(ctx)
		if err != nil {
			log.Fatalf("Failed to connect to SQLite database: %v", err)
		}
		log.Println("Using SQLite database for development")
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

	// Start gRPC server in a goroutine

	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen on gRPC port: %v", err)
		}

		fmt.Printf("gRPC server starting on :%s\n", grpcPort)
		fmt.Println("Services registered:")
		fmt.Println("  - ImageService")
		fmt.Println("  - LocationService")
		fmt.Println("  - gRPC Reflection enabled")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Wait for gRPC server to be ready
	fmt.Println("Waiting for gRPC server to start...")
	for i := 0; i < 30; i++ {
		if conn, err := net.Dial("tcp", "localhost:"+grpcPort); err == nil {
			conn.Close()
			fmt.Println("gRPC server is ready!")
			break
		}
		if i == 29 {
			log.Fatalf("gRPC server failed to start within 30 seconds")
		}
		time.Sleep(1 * time.Second)
	}

	// Create gRPC client connection for HTTP handler
	grpcConn, err := grpc.NewClient("localhost:"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer grpcConn.Close()

	// Create HTTP handler
	handler := handlers.NewHTTPHandler(grpcConn, grpcConn)

	// Setup routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Setup CORS middleware
	corsConfig := middleware.GetCORSConfig()

	// Get CORS origins from Secret Manager if available
	if secretManager != nil {
		corsOrigins := secretManager.GetSecretWithFallback("CORS_ALLOWED_ORIGINS", "CORS_ALLOWED_ORIGINS")
		if corsOrigins != "" {
			corsConfig.AllowedOrigins = []string{corsOrigins}
		}
	}

	// Wrap the mux with CORS middleware
	handlerWithCORS := middleware.CORS(corsConfig)(mux)

	// Start HTTP server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("HTTP server starting on :%s\n", port)
	fmt.Printf("Connected to gRPC server at: localhost:%s\n", grpcPort)
	fmt.Println("CORS enabled for origins:", corsConfig.AllowedOrigins)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /api/v1/images/current")
	fmt.Println("  POST /api/v1/images/upload")
	fmt.Println("  GET  /api/v1/images/count")
	fmt.Println("  GET  /api/v1/images")
	fmt.Println("  GET  /api/v1/images/{id}")
	fmt.Println("  DELETE /api/v1/images/{id}")
	fmt.Println("  GET  /api/v1/location/coords?lat=37.7749&lng=-122.4194")
	fmt.Println("  GET  /api/v1/location/name?name=San Francisco")
	fmt.Println("  GET  /health")

	// Clean up Secret Manager
	if secretManager != nil {
		defer secretManager.Close()
	}

	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}
