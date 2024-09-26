package services

import (
	"os"
	"path/filepath"
	"testing"
)

// Setup function to create the test environment for images
func setupImageTestEnvironment() error {
	return os.MkdirAll(uploadDir, os.ModePerm)
}

// TestProcessImageConcurrently tests the ProcessImageConcurrently function
func TestProcessImageConcurrently(t *testing.T) {
	// Setup test environment
	if err := setupImageTestEnvironment(); err != nil {
		t.Fatalf("Failed to set up test environment: %v", err)
	}

	t.Run("SuccessfulImageResize", func(t *testing.T) {
		// Create a result channel to receive the processed image
		resultCh := make(chan Result)

		// Get the test image path
		imagePath := filepath.Join("../testdata", "test.jpg")

		// Process the image (resize to 100x100)
		ProcessImageConcurrently(imagePath, 100, 100, resultCh)

		// Get the result
		result := <-resultCh

		// Check for errors
		if result.Err != nil {
			t.Fatalf("Expected no error, got: %v", result.Err)
		}

		// Check that the image was resized correctly
		if result.Img.Bounds().Dx() != 100 || result.Img.Bounds().Dy() != 100 {
			t.Fatalf("Expected image dimensions to be 100x100, got: %vx%v", result.Img.Bounds().Dx(), result.Img.Bounds().Dy())
		}
	})

	t.Run("ResizeWithOneDimension", func(t *testing.T) {
		// Create a result channel to receive the processed image
		resultCh := make(chan Result)

		// Get the test image path
		imagePath := filepath.Join("../testdata", "test.jpg")

		// Process the image (resize width to 100, height to 0 to preserve aspect ratio)
		ProcessImageConcurrently(imagePath, 100, 0, resultCh)

		// Get the result
		result := <-resultCh

		// Check for errors
		if result.Err != nil {
			t.Fatalf("Expected no error, got: %v", result.Err)
		}

		// Check that the width is resized to 100 and the height was adjusted proportionally
		if result.Img.Bounds().Dx() != 100 {
			t.Fatalf("Expected image width to be 100, got: %v", result.Img.Bounds().Dx())
		}
	})

	t.Run("InvalidImagePath", func(t *testing.T) {
		// Create a result channel to receive the processed image
		resultCh := make(chan Result)

		// Use an invalid image path
		imagePath := "invalid/path.jpg"

		// Process the image with an invalid path
		ProcessImageConcurrently(imagePath, 100, 100, resultCh)

		// Get the result
		result := <-resultCh

		// Check for error
		if result.Err == nil {
			t.Fatalf("Expected error for invalid image path, got nil")
		}
	})
}
