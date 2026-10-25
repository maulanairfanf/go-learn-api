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

var (
	ErrBlobTokenMissing = errors.New("BLOB_READ_WRITE_TOKEN is not configured")
	ErrBlobURLMissing   = errors.New("BLOB_URL is not configured")
)

type ImageStorage interface {
	Upload(ctx context.Context, filename, contentType string, size int64, body io.Reader) (string, error)
	Delete(ctx context.Context, blobURL string) error
}

type VercelBlobStorage struct {
	client *http.Client
}

func NewVercelBlobStorage() *VercelBlobStorage {
	return &VercelBlobStorage{client: http.DefaultClient}
}

func (s *VercelBlobStorage) Upload(ctx context.Context, filename, contentType string, size int64, body io.Reader) (string, error) {
	pathname, err := imagePath(filename)
	if err != nil {
		return "", err
	}

	baseURL, err := blobBaseURL()
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("pathname", pathname)
	request, err := s.newRequest(ctx, http.MethodPut, baseURL+"/?"+params.Encode(), body)
	if err != nil {
		return "", err
	}

	request.ContentLength = size
	request.Header.Set("x-vercel-blob-access", "public")
	request.Header.Set("x-add-random-suffix", "1")
	request.Header.Set("x-content-type", contentType)
	request.Header.Set("x-content-length", fmt.Sprintf("%d", size))

	response, err := s.client.Do(request)
	if err != nil {
		return "", err
	}

	responseBody, err := readResponseBody(response, "upload")
	if err != nil {
		return "", err
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", fmt.Errorf("decode vercel blob response: %w", err)
	}
	if result.URL == "" {
		return "", errors.New("vercel blob response did not contain a URL")
	}

	return result.URL, nil
}

func (s *VercelBlobStorage) Delete(ctx context.Context, blobURL string) error {
	body, err := json.Marshal(struct {
		URLs []string `json:"urls"`
	}{
		URLs: []string{blobURL},
	})
	if err != nil {
		return err
	}

	baseURL, err := blobBaseURL()
	if err != nil {
		return err
	}

	request, err := s.newRequest(ctx, http.MethodPost, baseURL+"/delete", bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := s.client.Do(request)
	if err != nil {
		return err
	}

	_, err = readResponseBody(response, "delete")
	return err
}

func (s *VercelBlobStorage) newRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	token := strings.TrimSpace(os.Getenv("BLOB_READ_WRITE_TOKEN"))
	if token == "" {
		return nil, ErrBlobTokenMissing
	}

	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("x-api-version", "13")
	if storeID := strings.TrimSpace(os.Getenv("BLOB_STORE_ID")); storeID != "" {
		request.Header.Set("x-vercel-blob-store-id", storeID)
	}

	return request, nil
}

func blobBaseURL() (string, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("BLOB_URL")), "/")
	if baseURL == "" {
		return "", ErrBlobURLMissing
	}
	return baseURL, nil
}

func readResponseBody(response *http.Response, operation string) ([]byte, error) {
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 4<<10))
	if err != nil {
		return nil, err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"vercel blob %s failed with status %d: %s",
			operation,
			response.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}

	return body, nil
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
