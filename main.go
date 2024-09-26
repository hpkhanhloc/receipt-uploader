package main

import (
	"log"
	"net/http"
	"os"
	"receipt-uploader/handlers"
	"receipt-uploader/models"
)

// Upload directory for receipts
const uploadDir = "uploads"

func main() {
	// Ensure the uploads directory exists
	os.MkdirAll(uploadDir, os.ModePerm)

	// Load receipts from JSON file into memory
	err := models.LoadReceiptsFromFile()
	if err != nil {
		log.Fatalf("Error loading receipts from file: %v", err)
	}

	// Define routes
	http.HandleFunc("/receipts", handlers.UploadReceipt) // POST /receipts to upload

	// Start server
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
