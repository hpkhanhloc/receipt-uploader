package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Custom error for invalid image uploads
var ErrInvalidImage = errors.New("not a valid image")

// SaveFile handles saving the uploaded file to the local filesystem
func SaveFile(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20) // Max 10MB
	if err != nil {
		log.Println("Error parsing multipart form:", err)
		return "", err
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file from request:", err)
		return "", err
	}
	defer file.Close()

	// Check if the file is an image by detecting its MIME type
	buffer := make([]byte, 512) // Buffer to store the first 512 bytes
	file.Read(buffer)           // Read the file into the buffer
	contentType := http.DetectContentType(buffer)

	if !strings.HasPrefix(contentType, "image/") {
		return "", ErrInvalidImage
	}

	// Rewind the file to the beginning for saving
	file.Seek(0, 0)

	fileID := uuid.New().String()
	filePath := filepath.Join("uploads", fileID+filepath.Ext(handler.Filename))

	// Create the file on the local filesystem
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

// GenerateReceiptID generates a unique receipt ID
func GenerateReceiptID() string {
	return uuid.New().String()
}
