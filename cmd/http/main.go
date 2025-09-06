package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/handlers"
	"github.com/NirvekPanda/Background-Image-Drive-API/internal/middleware"
	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get gRPC server address from environment or use default
	grpcAddr := os.Getenv("GRPC_SERVER_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	// Connect to gRPC services with retry mechanism
	var conn *grpc.ClientConn

	fmt.Printf("Attempting to connect to gRPC server at: %s\n", grpcAddr)

	for i := 0; i < 10; i++ {
		conn, err = grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			fmt.Printf("Successfully connected to gRPC server!\n")
			break
		}
		log.Printf("Failed to connect to gRPC server (attempt %d/10): %v", i+1, err)
		if i < 9 {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server after 10 attempts: %v\n"+
			"Make sure the gRPC server is running by executing: go run cmd/server/main.go", err)
	}
	defer conn.Close()

	// Test gRPC connection with a simple health check
	fmt.Println("Testing gRPC connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to call a simple gRPC method to verify connection
	imageClient := pb.NewImageServiceClient(conn)
	_, err = imageClient.GetImageCount(ctx, &pb.GetImageCountRequest{})
	if err != nil {
		log.Fatalf("gRPC connection test failed: %v\n"+
			"The gRPC server is not responding properly. Make sure it's running and healthy.", err)
	}
	fmt.Println("gRPC connection test successful!")

	// Create HTTP handler (both services use the same connection)
	handler := handlers.NewHTTPHandler(conn, conn)

	// Setup routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Setup CORS middleware
	corsConfig := middleware.GetCORSConfig()

	// Wrap the mux with CORS middleware
	handlerWithCORS := middleware.CORS(corsConfig)(mux)

	// Start HTTP server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("HTTP server starting on :%s\n", port)
	fmt.Printf("Connected to gRPC server at: %s\n", grpcAddr)
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

	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}
