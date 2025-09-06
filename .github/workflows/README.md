# 🚀 GitHub Actions Workflows

This directory contains GitHub Actions workflows for automated testing and deployment.

## 📋 Available Workflows

### 1. `ci.yml` - Continuous Integration
- **Triggers**: Push to main/develop, Pull requests
- **Purpose**: Runs unit tests, linting, security scanning, and builds
- **Jobs**: test, lint, build, security

### 2. `cd.yml` - Continuous Deployment  
- **Triggers**: Push to main branch
- **Purpose**: Builds Docker image and deploys to Google Cloud Run
- **Jobs**: deploy, notify

### 3. `test.yml` - API Integration Tests
- **Triggers**: Push to main/develop, Pull requests, Manual dispatch
- **Purpose**: Runs comprehensive API tests against production
- **Jobs**: api-tests
- **Tests**: Health check, Images upload/get, Location services

### 4. `manual-test.yml` - Manual API Testing
- **Triggers**: Manual dispatch only
- **Purpose**: Run specific API tests on demand
- **Options**: all, health, images, location

## 🧪 Running Tests

### Automatic Testing
Tests run automatically on:
- Push to `main` or `develop` branches
- Pull requests to `main` branch

### Manual Testing
1. Go to **Actions** tab in GitHub
2. Select **Manual API Test** workflow
3. Click **Run workflow**
4. Choose test type:
   - **all** - Run all tests
   - **health** - Health check only
   - **images** - Image upload/get only
   - **location** - Location services only

## 📊 Test Results

All workflows provide detailed output:
- ✅ **Success** - Test passed
- ❌ **Failure** - Test failed  
- ⚠️ **Warning** - Test passed with warnings

## 🔧 Configuration

Tests use environment variables from `testing/test.env`:
- `API_BASE_URL` - Production API endpoint
- `TEST_IMAGE_PATH` - Test image file
- `TEST_*` - Test data configuration

## 📈 Monitoring

Check workflow status in:
- **Actions** tab for detailed logs
- **Pull request** checks for PR status
- **Email notifications** for failures
