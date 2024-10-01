package services

import (
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// Custom error for invalid image uploads
var ErrInvalidImage = errors.New("not a valid image")

// SaveFile handles saving the uploaded file to the local filesystem
func SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Check if the file is an image by detecting its MIME type
	buffer := make([]byte, 512) // Buffer to store the first 512 bytes
	file.Read(buffer)           // Read the file into the buffer to detect content type
	contentType := http.DetectContentType(buffer)

	// Ensure the content type starts with "image/"
	if !strings.HasPrefix(contentType, "image/") {
		return "", ErrInvalidImage
	}

	// Rewind the file after reading its MIME type
	file.Seek(0, 0)

	// Create the file on the filesystem
	fileID := GenerateReceiptID()
	filePath := filepath.Join("uploads", fileID+filepath.Ext(fileHeader.Filename))
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return "", fmt.Errorf("failed to create file on the server: %v", err)
	}
	defer f.Close()

	// Copy the uploaded file to the filesystem
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println("Error copying file to filesystem:", err)
		return "", fmt.Errorf("failed to copy file to the server: %v", err)
	}

	log.Println("File saved successfully:", filePath)
	return filePath, nil
}

// SaveImage saves the resized image to the specified file path
func SaveImage(img image.Image, filePath string) error {
	// Create the file on the filesystem
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Save the image as a JPEG file
	err = imaging.Save(img, filePath)
	if err != nil {
		return err
	}

	return nil
}

// GenerateReceiptID generates a unique receipt ID
func GenerateReceiptID() string {
	return uuid.New().String()
}
