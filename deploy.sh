#!/bin/bash

# Deployment script for Google Cloud Run with Secret Manager
set -e

# Configuration
PROJECT_ID="portfolio-420-69"
SERVICE_NAME="background-image-drive-api"
REGION="us-west1"
IMAGE_NAME="gcr.io/$PROJECT_ID/$SERVICE_NAME"

echo "🚀 Deploying $SERVICE_NAME to Google Cloud Run..."

# Set the project
gcloud config set project $PROJECT_ID

# Build and push Docker image
echo "📦 Building and pushing Docker image..."
docker build --platform linux/amd64 -t $IMAGE_NAME .
docker push $IMAGE_NAME

# Deploy to Cloud Run
echo "🚀 Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image $IMAGE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --add-cloudsql-instances=$PROJECT_ID:us-central1:portfolio-images-db \
  --memory=1Gi \
  --cpu=1 \
  --max-instances=10 \
  --min-instances=0

echo "✅ Deployment complete!"
echo "🌐 Service URL: https://$SERVICE_NAME-189526192204.$REGION.run.app"
echo ""
echo "📋 Next steps:"
echo "1. Make sure your service account has Secret Manager access:"
echo "   gcloud projects add-iam-policy-binding $PROJECT_ID \\"
echo "     --member=\"serviceAccount:YOUR_SERVICE_ACCOUNT@$PROJECT_ID.iam.gserviceaccount.com\" \\"
echo "     --role=\"roles/secretmanager.secretAccessor\""
echo ""
echo "2. Test the health endpoint:"
echo "   curl https://$SERVICE_NAME-189526192204.$REGION.run.app/health"