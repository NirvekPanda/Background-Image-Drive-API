# 🚀 Portfolio Images API - Cloud Run Deployment

## 📋 Prerequisites

1. **Google Cloud CLI** installed and authenticated
2. **Docker** installed locally
3. **Environment variables** configured in `.env` file
4. **OAuth2 credentials** for Google Drive access

## 🔧 Environment Variables

Make sure your `.env` file contains:

```bash
# Google Drive Configuration
GOOGLE_DRIVE_FOLDER_ID=your-folder-id-here

# Google Maps API
GOOGLE_MAPS_API_KEY=your-maps-api-key

# Cloud SQL Configuration
CLOUD_SQL_CONNECTION_NAME=your-cloud-sql-ip
CLOUD_SQL_DATABASE=portfolio_images
CLOUD_SQL_USER=portfolio_user
CLOUD_SQL_PASSWORD=your-password

# OAuth2 Credentials (create these files)
# oauth_credentials.json - Download from Google Cloud Console
# token.json - Generated after first OAuth2 authentication
```

## 🐳 Local Docker Testing

### 1. Test Docker Container Locally

```bash
# Load environment variables
source .env

# Run local Docker test
./test-docker.sh
```

### 2. Manual Docker Testing

```bash
# Build image
docker build -t portfolio-images-api .

# Run container
docker run -d \
  --name portfolio-images-test \
  -p 8080:8080 \
  -p 50051:50051 \
  -e GOOGLE_DRIVE_FOLDER_ID="$GOOGLE_DRIVE_FOLDER_ID" \
  -e GOOGLE_MAPS_API_KEY="$GOOGLE_MAPS_API_KEY" \
  -e CLOUD_SQL_CONNECTION_NAME="$CLOUD_SQL_CONNECTION_NAME" \
  -e CLOUD_SQL_DATABASE="$CLOUD_SQL_DATABASE" \
  -e CLOUD_SQL_USER="$CLOUD_SQL_USER" \
  -e CLOUD_SQL_PASSWORD="$CLOUD_SQL_PASSWORD" \
  portfolio-images-api

# Test API
curl http://localhost:8080/health
```

## ☁️ Cloud Run Deployment

### 1. Deploy to Cloud Run

```bash
# Load environment variables
source .env

# Deploy to Cloud Run
./deploy.sh
```

### 2. Manual Cloud Run Deployment

```bash
# Set project
gcloud config set project portfolio-420-69

# Build and push image
docker build -t gcr.io/portfolio-420-69/portfolio-images-api .
docker push gcr.io/portfolio-420-69/portfolio-images-api

# Deploy to Cloud Run
gcloud run deploy portfolio-images-api \
  --image gcr.io/portfolio-420-69/portfolio-images-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --max-instances 10 \
  --min-instances 0 \
  --timeout 300 \
  --concurrency 80 \
  --set-env-vars "GOOGLE_DRIVE_FOLDER_ID=$GOOGLE_DRIVE_FOLDER_ID" \
  --set-env-vars "GOOGLE_MAPS_API_KEY=$GOOGLE_MAPS_API_KEY" \
  --set-env-vars "CLOUD_SQL_CONNECTION_NAME=$CLOUD_SQL_CONNECTION_NAME" \
  --set-env-vars "CLOUD_SQL_DATABASE=$CLOUD_SQL_DATABASE" \
  --set-env-vars "CLOUD_SQL_USER=$CLOUD_SQL_USER" \
  --set-env-vars "CLOUD_SQL_PASSWORD=$CLOUD_SQL_PASSWORD"
```

## 🔐 OAuth2 Setup for Cloud Run

### 1. Create OAuth2 Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Navigate to: APIs & Services → Credentials
3. Create OAuth 2.0 Client ID
4. Application type: Web application
5. Authorized redirect URIs: `http://localhost:8080/oauth/callback`
6. Download JSON and save as `oauth_credentials.json`

### 2. Generate Initial Token

```bash
# Run locally to generate token
go run cmd/server/main.go

# Follow OAuth2 flow to generate token.json
# Copy token.json to your project root
```

### 3. Secure Credential Storage

**⚠️ IMPORTANT: Never commit OAuth2 credentials to version control!**

The following files are automatically ignored by `.gitignore`:
- `oauth_credentials.json` - OAuth2 client credentials
- `token.json` - OAuth2 access token
- `.env` - Environment variables

For production deployment, credentials are stored securely in Google Cloud Secret Manager.

## 📊 Monitoring and Logs

### View Logs

```bash
# Cloud Run logs
gcloud logs read --service=portfolio-images-api --limit=50

# Real-time logs
gcloud logs tail --service=portfolio-images-api
```

### Monitor Performance

```bash
# Service details
gcloud run services describe portfolio-images-api --region=us-central1

# Service URL
gcloud run services describe portfolio-images-api --region=us-central1 --format='value(status.url)'
```

## 🧪 Testing Deployed API

```bash
# Get service URL
SERVICE_URL=$(gcloud run services describe portfolio-images-api --region=us-central1 --format='value(status.url)')

# Test health
curl $SERVICE_URL/health

# Test image count
curl $SERVICE_URL/api/v1/images/count

# Test image upload
curl -X POST $SERVICE_URL/api/v1/images/upload \
  -F "title=Test Upload" \
  -F "description=Testing Cloud Run" \
  -F "image=@amsterdam.jpg"
```

## 🔧 Troubleshooting

### Common Issues

1. **OAuth2 Token Expired**
   - Regenerate token locally
   - Update Cloud Run with new token

2. **Database Connection Failed**
   - Check Cloud SQL instance is running
   - Verify IP address and credentials

3. **Google Drive Upload Failed**
   - Verify OAuth2 credentials
   - Check folder permissions

### Debug Commands

```bash
# Check container logs
docker logs portfolio-images-test

# Check Cloud Run logs
gcloud logs read --service=portfolio-images-api --limit=100

# Test database connection
gcloud sql connect portfolio-images-db --user=portfolio_user --database=portfolio_images
```

## 📈 Scaling

### Auto-scaling Configuration

- **Min instances**: 0 (cost-effective)
- **Max instances**: 10 (handle traffic spikes)
- **Memory**: 1Gi (sufficient for image processing)
- **CPU**: 1 (efficient for API operations)
- **Concurrency**: 80 (optimal for Cloud Run)

### Cost Optimization

- Use min instances = 0 for development
- Set max instances based on expected traffic
- Monitor usage with Cloud Monitoring
- Use Cloud SQL free tier (db-f1-micro)
