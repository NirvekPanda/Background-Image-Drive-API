package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	pb "github.com/NirvekPanda/Background-Image-Drive-API/proto/gen"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// DriveUtilOAuth handles Google Drive operations using OAuth2
type DriveUtilOAuth struct {
	service  *drive.Service
	folderID string // ID of the specific folder to use
}

// OAuthConfig holds OAuth2 configuration
type OAuthConfig struct {
	Web struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
	} `json:"web"`
}

// NewDriveUtilOAuth creates a new Google Drive utility client using OAuth2
func NewDriveUtilOAuth(ctx context.Context, oauthConfigPath string, tokenPath string, folderID string) (*DriveUtilOAuth, error) {
	configData, err := os.ReadFile(oauthConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OAuth config: %v", err)
	}

	var oauthConfig OAuthConfig
	if err := json.Unmarshal(configData, &oauthConfig); err != nil {
		return nil, fmt.Errorf("failed to parse OAuth config: %v", err)
	}

	config := &oauth2.Config{
		ClientID:     oauthConfig.Web.ClientID,
		ClientSecret: oauthConfig.Web.ClientSecret,
		RedirectURL:  oauthConfig.Web.RedirectURIs[0],
		Scopes: []string{
			drive.DriveFileScope,
		},
		Endpoint: google.Endpoint,
	}

	token, err := loadOrGetToken(ctx, config, tokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	client := config.Client(ctx, token)
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Drive service: %v", err)
	}

	return &DriveUtilOAuth{
		service:  service,
		folderID: folderID,
	}, nil
}

func loadOrGetToken(ctx context.Context, config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	if tokenData, err := os.ReadFile(tokenPath); err == nil {
		var token oauth2.Token
		if err := json.Unmarshal(tokenData, &token); err == nil {
			if token.Valid() {
				return &token, nil
			}
			if token, err := config.TokenSource(ctx, &token).Token(); err == nil {
				_ = saveToken(token, tokenPath)
				return token, nil
			}
		}
	}

	return startOAuthFlow(ctx, config, tokenPath)
}

func startOAuthFlow(ctx context.Context, config *oauth2.Config, tokenPath string) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("OAuth2 authentication required. Please visit: %s\n", authURL)
	fmt.Print("Enter authorization code: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("failed to read authorization code: %v", err)
	}

	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	if err := saveToken(token, tokenPath); err != nil {
		return nil, fmt.Errorf("failed to save token: %v", err)
	}

	return token, nil
}

func saveToken(token *oauth2.Token, tokenPath string) error {
	tokenData, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(tokenPath, tokenData, 0600)
}

func (d *DriveUtilOAuth) UploadFile(ctx context.Context, filename string, imageData []byte) (string, error) {
	file := &drive.File{
		Name:    filename,
		Parents: []string{d.folderID},
	}

	call := d.service.Files.Create(file).Media(bytes.NewReader(imageData)).Context(ctx)
	createdFile, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	return createdFile.Id, nil
}

func (d *DriveUtilOAuth) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	call := d.service.Files.Get(fileID).Context(ctx)
	resp, err := call.Download()
	if err != nil {
		return nil, fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %v", err)
	}

	return data, nil
}

func (d *DriveUtilOAuth) DeleteFile(ctx context.Context, fileID string) error {
	err := d.service.Files.Delete(fileID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("delete failed: %v", err)
	}
	return nil
}

func (d *DriveUtilOAuth) ListFilesInFolder(ctx context.Context) ([]*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents and trashed=false", d.folderID)

	call := d.service.Files.List().
		Q(query).
		Fields("files(id,name,mimeType,createdTime,modifiedTime,size)").
		Context(ctx)

	files, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	return files.Files, nil
}

func (d *DriveUtilOAuth) ListImageFilesInFolder(ctx context.Context) ([]*pb.ImageMetadata, error) {
	query := fmt.Sprintf("'%s' in parents and trashed=false and (mimeType contains 'image/')", d.folderID)

	call := d.service.Files.List().
		Q(query).
		Fields("files(id,name,mimeType,createdTime,modifiedTime,size)").
		Context(ctx)

	files, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list image files: %v", err)
	}

	var imageMetadata []*pb.ImageMetadata
	for _, file := range files.Files {
		title := file.Name
		if lastDot := strings.LastIndex(file.Name, "."); lastDot != -1 {
			title = file.Name[:lastDot]
		}

		metadata := &pb.ImageMetadata{
			Id:          file.Id,
			Title:       title,
			Description: fmt.Sprintf("Image uploaded on %s", file.CreatedTime),
			DriveFileId: file.Id,
			Location: &pb.Location{
				Name: "Unknown Location",
			},
		}

		imageMetadata = append(imageMetadata, metadata)
	}

	return imageMetadata, nil
}

func (d *DriveUtilOAuth) GetFolderInfo(ctx context.Context) (*drive.File, error) {
	call := d.service.Files.Get(d.folderID).
		Fields("id,name,createdTime,modifiedTime").
		Context(ctx)

	folder, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get folder info: %v", err)
	}

	return folder, nil
}
