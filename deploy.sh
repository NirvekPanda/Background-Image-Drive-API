#!/bin/bash

# Configuration
PROJECT_ID="portfolio-420-69"
SERVICE_NAME="background-image-drive-api"
REGION="us-west1"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "🚀 Deploying Portfolio Images API to Cloud Run"
echo "Project: $PROJECT_ID"
echo "Service: $SERVICE_NAME"
echo "Region: $REGION"
echo "Image: $IMAGE_NAME"
echo ""

# Build and push Docker image
echo "📦 Building Docker image..."
docker build -t $IMAGE_NAME .

echo "📤 Pushing image to Google Container Registry..."
docker push $IMAGE_NAME

# Create secrets in Google Cloud Secret Manager (if they don't exist)
echo "🔐 Setting up secrets in Google Cloud Secret Manager..."
gcloud secrets create oauth-credentials --data-file=oauth_credentials.json --quiet || echo "Secret oauth-credentials already exists"
gcloud secrets create oauth-token --data-file=token.json --quiet || echo "Secret oauth-token already exists"

# Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 1Gi \
  --cpu 1 \
  --max-instances 10 \
  --min-instances 0 \
  --timeout 300 \
  --concurrency 80 \
  --env-vars-file env-vars.yaml \
  --set-secrets "oauth_credentials.json=oauth-credentials:latest" \
  --set-secrets "token.json=oauth-token:latest"

echo ""
echo "✅ Deployment complete!"
echo "🌐 Your API is now available at:"
gcloud run services describe $SERVICE_NAME --region $REGION --format 'value(status.url)'
