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
	http.HandleFunc("/receipts", handleReceipts)       // unified route for both POST and GET methods on /receipts
	http.HandleFunc("/receipts/", handlers.GetReceipt) // GET /receipts/{receipt_id} to retrieve

	// Start server
	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// handleReceipts handles both POST (upload) and GET (list receipts) methods on /receipts
func handleReceipts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlers.UploadReceipt(w, r)
	case http.MethodGet:
		handlers.ListReceipts(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
