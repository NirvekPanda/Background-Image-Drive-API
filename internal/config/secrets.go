package config

import (
	"context"
	"fmt"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// SecretManager handles Google Cloud Secret Manager operations
type SecretManager struct {
	client *secretmanager.Client
	ctx    context.Context
}

// NewSecretManager creates a new SecretManager instance
func NewSecretManager(ctx context.Context) (*SecretManager, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret manager client: %v", err)
	}

	return &SecretManager{
		client: client,
		ctx:    ctx,
	}, nil
}

// GetSecret retrieves a secret from Google Cloud Secret Manager
func (sm *SecretManager) GetSecret(secretName string) (string, error) {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		// Try to get project ID from other environment variables
		projectID = os.Getenv("GCP_PROJECT")
		if projectID == "" {
			projectID = os.Getenv("GCLOUD_PROJECT")
			if projectID == "" {
				// In Cloud Run, the project ID is available in the metadata
				// For now, hardcode it since we know it
				projectID = "portfolio-420-69"
			}
		}
	}

	// Create the request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", projectID, secretName),
	}

	// Call the API
	result, err := sm.client.AccessSecretVersion(sm.ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to access secret %s: %v", secretName, err)
	}

	return string(result.Payload.Data), nil
}

// GetSecretWithFallback retrieves a secret from Secret Manager, falls back to environment variable
func (sm *SecretManager) GetSecretWithFallback(secretName, envVar string) string {
	// Try to get from Secret Manager first
	secret, err := sm.GetSecret(secretName)
	if err == nil {
		return secret
	}

	// Fall back to environment variable
	return os.Getenv(envVar)
}

// Close closes the Secret Manager client
func (sm *SecretManager) Close() error {
	return sm.client.Close()
}
