package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NirvekPanda/Background-Image-Drive-API/internal/handlers"
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

	// Connect to gRPC services (both services run on the same port)
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Create HTTP handler (both services use the same connection)
	handler := handlers.NewHTTPHandler(conn, conn)

	// Setup routes
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	// Start HTTP server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("HTTP server starting on :%s\n", port)
	fmt.Printf("Connected to gRPC server at: %s\n", grpcAddr)
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

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
