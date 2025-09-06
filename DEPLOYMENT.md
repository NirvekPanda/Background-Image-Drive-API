# Deployment Guide

## Environment Variables Setup

### 1. Create Environment File
Copy the template and fill in your values:
```bash
cp env-vars.yaml.template env-vars.yaml
```

### 2. Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed origins | `"https://yourdomain.com,http://localhost:3000"` |
| `ENVIRONMENT` | Environment name | `"production"` |
| `GOOGLE_DRIVE_FOLDER_ID` | Google Drive folder ID for storing images | `"1GFVR5bXu5JOtVwNnpY8JrnnPRE7XKzYQ"` |
| `GOOGLE_MAPS_API_KEY` | Google Maps API key | `"AIzaSyBYTQ-2qwwVZabGQpMF2fZNtc9SmtbCbgs"` |
| `CLOUD_SQL_CONNECTION_NAME` | Cloud SQL connection name | `"project-id:region:instance-name"` |
| `CLOUD_SQL_DATABASE` | Database name | `"portfolio_images"` |
| `CLOUD_SQL_USER` | Database username | `"portfolio_user"` |
| `CLOUD_SQL_PASSWORD` | Database password | `"your-secure-password"` |
| `GRPC_PORT` | Internal gRPC port | `"50051"` |

### 3. Google Cloud Credentials

Place your OAuth credentials in the appropriate location:
- **Local development**: `oauth_credentials.json` and `token.json` in project root
- **Production**: Mount as secrets in `/app/secrets/`

## Local Development

1. **Set up environment variables:**
   ```bash
   cp env-vars.yaml.template env-vars.yaml
   # Edit env-vars.yaml with your values
   ```

2. **Add Google OAuth credentials:**
   - Download `oauth_credentials.json` from Google Cloud Console
   - Place in project root
   - Run the app once to generate `token.json`

3. **Run the service:**
   ```bash
   go run cmd/combined/main.go
   ```

## Google Cloud Run Deployment

### 1. Create Secrets in Google Cloud Secret Manager
```bash
# Set your project ID
export PROJECT_ID=your-project-id

# Create all secrets
gcloud secrets create CORS_ALLOWED_ORIGINS --data-file=- <<< "https://yourdomain.com,http://localhost:3000"
gcloud secrets create ENVIRONMENT --data-file=- <<< "production"
gcloud secrets create GOOGLE_DRIVE_FOLDER_ID --data-file=- <<< "your-folder-id"
gcloud secrets create GOOGLE_MAPS_API_KEY --data-file=- <<< "your-maps-key"
gcloud secrets create CLOUD_SQL_CONNECTION_NAME --data-file=- <<< "project-id:region:instance-name"
gcloud secrets create CLOUD_SQL_DATABASE --data-file=- <<< "portfolio_images"
gcloud secrets create CLOUD_SQL_USER --data-file=- <<< "portfolio_user"
gcloud secrets create CLOUD_SQL_PASSWORD --data-file=- <<< "your-db-password"
gcloud secrets create GRPC_PORT --data-file=- <<< "50051"
```

### 2. Build and Push Docker Image
```bash
# Build the image
docker build -t gcr.io/$PROJECT_ID/portfolio-images .

# Push to Google Container Registry
docker push gcr.io/$PROJECT_ID/portfolio-images
```

### 3. Deploy to Cloud Run
```bash
gcloud run deploy portfolio-images \
  --image gcr.io/$PROJECT_ID/portfolio-images \
  --platform managed \
  --region us-west1 \
  --allow-unauthenticated \
  --add-cloudsql-instances=project-id:region:instance-name \
  --service-account=your-service-account@project-id.iam.gserviceaccount.com
```

### 4. Grant Secret Manager Access
```bash
# Grant the service account access to Secret Manager
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:your-service-account@project-id.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"
```

### 5. Set up Cloud SQL Connection
Make sure your Cloud Run service has access to Cloud SQL:
```bash
gcloud run services add-iam-policy-binding portfolio-images \
  --member="serviceAccount:your-service-account@project-id.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"
```

## Security Notes

- **Never commit** `env-vars.yaml` to version control
- **Never include** API keys or passwords in Docker images
- Use Google Cloud Secret Manager for production secrets
- Rotate API keys and passwords regularly
- Enable Cloud SQL private IP when possible

## Troubleshooting

### Common Issues

1. **Cloud SQL Connection Failed**
   - Verify connection name format: `project-id:region:instance-name`
   - Check Cloud SQL instance is running
   - Ensure Cloud Run has Cloud SQL access

2. **Google Drive API Errors**
   - Verify OAuth credentials are correct
   - Check Google Drive folder ID exists
   - Ensure proper API permissions

3. **CORS Errors**
   - Add your frontend domain to `CORS_ALLOWED_ORIGINS`
   - Check environment variable is set correctly