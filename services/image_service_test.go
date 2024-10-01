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

// TestProcessImage tests the ProcessImage function
func TestProcessImage(t *testing.T) {
	// Setup test environment
	if err := setupImageTestEnvironment(); err != nil {
		t.Fatalf("Failed to set up test environment: %v", err)
	}

	t.Run("SuccessfulImageResize", func(t *testing.T) {
		// Get the test image path
		imagePath := filepath.Join("../testdata", "test.jpg")

		// Process the image (resize to 100x100)
		img, err := ProcessImage(imagePath, 100, 100)

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check that the image was resized correctly
		if img.Bounds().Dx() != 100 || img.Bounds().Dy() != 100 {
			t.Fatalf("Expected image dimensions to be 100x100, got: %vx%v", img.Bounds().Dx(), img.Bounds().Dy())
		}
	})

	t.Run("ResizeWithOneDimension", func(t *testing.T) {
		// Get the test image path
		imagePath := filepath.Join("../testdata", "test.jpg")

		// Process the image (resize width to 100, height to 0 to preserve aspect ratio)
		img, err := ProcessImage(imagePath, 100, 0)

		// Check for errors
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Check that the width is resized to 100 and the height was adjusted proportionally
		if img.Bounds().Dx() != 100 {
			t.Fatalf("Expected image width to be 100, got: %v", img.Bounds().Dx())
		}
		if img.Bounds().Dy() == 0 {
			t.Fatalf("Expected image height is not 0, got: %v", img.Bounds().Dx())
		}
	})

	t.Run("InvalidImagePath", func(t *testing.T) {
		// Use an invalid image path
		imagePath := "invalid/path.jpg"

		// Process the image with an invalid path
		_, err := ProcessImage(imagePath, 100, 100)

		// Check for error
		if err == nil {
			t.Fatalf("Expected error for invalid image path, got nil")
		}
	})
}
