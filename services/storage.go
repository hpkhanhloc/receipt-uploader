package services

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

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

	fileID := uuid.New().String()
	filePath := filepath.Join("uploads", fileID+filepath.Ext(handler.Filename))

	// Create the file on the local filesystem
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating file:", err)
		return "", err
	}
	defer f.Close()

	// Copy the uploaded file to the filesystem
	_, err = io.Copy(f, file)
	if err != nil {
		log.Println("Error copying file to filesystem:", err)
		return "", err
	}

	log.Println("File saved successfully:", filePath)
	return filePath, nil
}

// GenerateReceiptID generates a unique receipt ID
func GenerateReceiptID() string {
	return uuid.New().String()
}
