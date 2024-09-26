package services

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const uploadDir = "../uploads"

// Setup function to create the uploads directory before tests
func setupTestEnvironment() error {
	return os.MkdirAll(uploadDir, os.ModePerm)
}

// Helper function to create a multipart request for file upload from a test file
func createMultipartRequest(fileName string) (*http.Request, error) {
	// Open the file from the testdata directory
	filePath := filepath.Join("../testdata/", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	// Copy the file contents into the form
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	writer.Close()

	// Create the HTTP request with the multipart body
	req := httptest.NewRequest(http.MethodPost, "/receipts", body)
	req.Header.Set("Content-Type", writer.FormDataContentType()) // Set Content-Type for multipart form

	return req, nil
}

// TestSaveFile tests the SaveFile function
func TestSaveFile(t *testing.T) {
	// Setup the test environment
	if err := setupTestEnvironment(); err != nil {
		t.Fatalf("Failed to create uploads directory: %v", err)
	}

	// Valid image test case
	t.Run("ValidImageUpload", func(t *testing.T) {
		req, err := createMultipartRequest("test.jpg")
		if err != nil {
			t.Fatalf("Failed to create multipart request: %v", err)
		}

		// Run the function
		filePath, err := SaveFile(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Defer cleanup in case of any errors
		defer os.Remove(filePath)

		// Verify the file was saved correctly
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("Expected file to be saved at %s, but it wasn't", filePath)
		}
	})

	// Non-image file test case
	t.Run("NonImageFileUpload", func(t *testing.T) {
		req, err := createMultipartRequest("test.txt")
		if err != nil {
			t.Fatalf("Failed to create multipart request: %v", err)
		}

		// Run the function and check for invalid image error
		_, err = SaveFile(req)
		if err == nil || !errors.Is(err, ErrInvalidImage) {
			t.Fatalf("Expected error for invalid image, got %v", err)
		}
	})
}
