# 🧪 Simple Test Suite

A clean, simple test suite for the Portfolio Images API.

## 📁 Test Files

### `test-health.sh`
Tests if the API is live and responsive.
- Health endpoint check
- Response time measurement
- CORS headers verification

### `test-images.sh`
Tests image upload and retrieval functionality.
- Upload image with metadata
- Get image by ID
- List all images
- Get image count

### `test-location.sh`
Tests location services and metadata extraction.
- Get location from coordinates
- Get location from name
- Test multiple location names

### `run-all-tests.sh`
Master test runner that executes all tests in sequence.

## 🚀 Quick Start

```bash
# Run all tests
./testing/run-all-tests.sh

# Run individual tests
./testing/test-health.sh
./testing/test-images.sh
./testing/test-location.sh
```

## ⚙️ Configuration

Set environment variables or create `testing/test.env`:

```bash
# API Configuration
API_BASE_URL=https://background-image-drive-api-189526192204.us-west1.run.app

# Test Image
TEST_IMAGE_PATH=testing/test_image.jpg

# Test Data
TEST_IMAGE_TITLE="My Test Image"
TEST_IMAGE_DESCRIPTION="Test upload description"
TEST_LATITUDE=52.3676
TEST_LONGITUDE=4.9041
TEST_LOCATION_NAME="Amsterdam"
TEST_COUNTRY="Netherlands"
TEST_CITY="Amsterdam"
TEST_LOCATION_NAMES="Amsterdam,London,Paris,Tokyo,New York"
```

## 📊 Test Results

All tests provide clear output:
- ✅ **Success** - Test passed
- ❌ **Failure** - Test failed
- ⚠️ **Warning** - Test passed with warnings

## 🔧 Features

- **Simple & Clean**: Easy to understand and maintain
- **Production Ready**: Tests against live API
- **Comprehensive**: Covers all major endpoints
- **Robust**: Handles different response formats
- **Informative**: Clear success/failure reporting