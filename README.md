# Portfolio Image Management API

A Go-based microservices API for managing portfolio images with Google Drive integration and location metadata.

## Features

- **Image Management**: Upload, retrieve, and delete images with metadata
- **Google Drive Integration**: Store images in Google Drive
- **Location Services**: Extract location data from coordinates or names
- **RESTful API**: HTTP endpoints for easy integration
- **gRPC Services**: High-performance internal communication
- **Metadata Storage**: Track image titles, descriptions, and locations

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Gateway  │────│   gRPC Server   │────│  Google Drive   │
│   (Port 8080)   │    │   (Port 50051)  │    │   Integration   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         │              ┌─────────────────┐
         └──────────────│  Location API   │
                        │   (Simplified)  │
                        └─────────────────┘
```

## Services

### ImageService
- `GetCurrentImage` - Get the most recently uploaded image
- `UploadImage` - Upload new images with metadata
- `GetImageCount` - Get total number of images
- `GetImageById` - Retrieve specific image by ID
- `DeleteImage` - Remove images from storage and Google Drive

### LocationService
- `GetLocationFromCoords` - Convert coordinates to location data
- `GetLocationFromName` - Convert location name to structured data

## API Endpoints

### Images
- `GET /api/v1/images/current` - Get current image
- `POST /api/v1/images/upload` - Upload new image
- `GET /api/v1/images/count` - Get image count
- `GET /api/v1/images/{id}` - Get image by ID
- `DELETE /api/v1/images/{id}` - Delete image

### Location
- `GET /api/v1/location/coords?lat=37.7749&lng=-122.4194` - Get location from coordinates
- `GET /api/v1/location/name?name=San Francisco` - Get location from name

### Health
- `GET /health` - Health check

## Setup

### Prerequisites

1. Go 1.24 or later
2. Google Cloud Project with Drive API enabled
3. Service Account credentials
4. Google Maps API key (for geocoding)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd portfolio-images
```

2. Install dependencies:
```bash
go mod tidy
```

3. Generate protobuf files:
```bash
make proto
```

4. Set up Google Drive credentials:
   - Create a service account in Google Cloud Console
   - Enable Google Drive API
   - Download the service account JSON key
   - Copy it to `credentials.json` or set `GOOGLE_DRIVE_CREDENTIALS_PATH`

5. Set up Google Maps API:
   - Go to Google Cloud Console → APIs & Services → Library
   - Enable "Geocoding API"
   - Go to APIs & Services → Credentials
   - Create an API key
   - Set `GOOGLE_MAPS_API_KEY` environment variable

### Running the Services

1. Start the gRPC server:
```bash
go run cmd/server/main.go
```

2. Start the HTTP gateway (in another terminal):
```bash
go run cmd/http/main.go
```

### Environment Variables

- `GRPC_PORT` - gRPC server port (default: 50051)
- `HTTP_PORT` - HTTP server port (default: 8080)
- `GOOGLE_DRIVE_CREDENTIALS_PATH` - Path to Google Drive credentials (default: credentials.json)
- `GOOGLE_MAPS_API_KEY` - Google Maps API key for geocoding (required)
- `GRPC_SERVER_ADDR` - gRPC server address for HTTP gateway (default: localhost:50051)

## Usage Examples

### Upload an Image
```bash
curl -X POST http://localhost:8080/api/v1/images/upload \
  -F "image=@photo.jpg" \
  -F "title=My Photo" \
  -F "description=A beautiful sunset" \
  -F "latitude=37.7749" \
  -F "longitude=-122.4194" \
  -F "location_name=San Francisco"
```

### Get Current Image
```bash
curl http://localhost:8080/api/v1/images/current
```

### Get Location from Coordinates
```bash
curl "http://localhost:8080/api/v1/location/coords?lat=37.7749&lng=-122.4194"
```

## Development

### Project Structure
```
├── cmd/
│   ├── server/          # gRPC server
│   └── http/            # HTTP gateway
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   └── services/        # Business logic
├── proto/               # Protocol buffer definitions
└── bin/                 # Built binaries
```

### Building
```bash
make build
```

### Testing
```bash
make test
```

## Future Enhancements

- [x] Database integration (PostgreSQL/MongoDB)
- [ ] Real geocoding service integration
- [ ] Image processing and resizing
- [ ] Authentication and authorization
- [ ] Docker containerization
- [ ] Kubernetes deployment
- [ ] Monitoring and logging
- [ ] Rate limiting
- [ ] Caching layer

