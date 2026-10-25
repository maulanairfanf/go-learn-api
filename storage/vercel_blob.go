package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var ErrBlobTokenMissing = errors.New("BLOB_READ_WRITE_TOKEN is not configured")

type ImageStorage interface {
	Upload(ctx context.Context, filename, contentType string, size int64, body io.Reader) (string, error)
	Delete(ctx context.Context, blobURL string) error
}

type VercelBlobStorage struct {
	client *http.Client
}

func NewVercelBlobStorage() *VercelBlobStorage {
	return &VercelBlobStorage{
		client: http.DefaultClient,
	}
}

func (s *VercelBlobStorage) Upload(ctx context.Context, filename, contentType string, size int64, body io.Reader) (string, error) {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	storeID := os.Getenv("BLOB_STORE_ID")
	if strings.TrimSpace(token) == "" {
		return "", ErrBlobTokenMissing
	}

	pathname, err := imagePath(filename)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("pathname", pathname)
	uploadURL := "https://vercel.com/api/blob/?" + params.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, body)
	if err != nil {
		return "", err
	}

	request.ContentLength = size
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("x-api-version", "13")
	request.Header.Set("x-vercel-blob-access", "public")
	request.Header.Set("x-add-random-suffix", "1")
	request.Header.Set("x-content-type", contentType)
	request.Header.Set("x-content-length", fmt.Sprintf("%d", size))
	if storeID != "" {
		request.Header.Set("x-vercel-blob-store-id", storeID)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 4<<10))
		return "", fmt.Errorf("vercel blob upload failed with status %d: %s", response.StatusCode, strings.TrimSpace(string(message)))
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode vercel blob response: %w", err)
	}
	if result.URL == "" {
		return "", errors.New("vercel blob response did not contain a URL")
	}

	return result.URL, nil
}

func (s *VercelBlobStorage) Delete(ctx context.Context, blobURL string) error {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if strings.TrimSpace(token) == "" {
		return ErrBlobTokenMissing
	}

	body, err := json.Marshal(struct {
		URLs []string `json:"urls"`
	}{
		URLs: []string{blobURL},
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://vercel.com/api/blob/delete",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-version", "13")

	if storeID := os.Getenv("BLOB_STORE_ID"); storeID != "" {
		request.Header.Set("x-vercel-blob-store-id", storeID)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 4<<10))
		return fmt.Errorf(
			"vercel blob delete failed with status %d: %s",
			response.StatusCode,
			strings.TrimSpace(string(message)),
		)
	}

	return nil
}

func imagePath(filename string) (string, error) {
	ext := filepath.Ext(filepath.Base(filename))
	if ext == "" {
		return "", errors.New("image filename must have an extension")
	}

	contentType := mime.TypeByExtension(strings.ToLower(ext))
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("unsupported image extension %q", ext)
	}

	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate image name: %w", err)
	}

	return "products/" + hex.EncodeToString(randomBytes) + ext, nil
}
