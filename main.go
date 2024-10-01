package main

import (
	"log"
	"net/http"
	"os"
	"receipt-uploader/handlers"
	"receipt-uploader/models"
	"strings"
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
	http.HandleFunc("/receipts", handleReceipts)         // unified route for both POST and GET methods on /receipts
	http.HandleFunc("/receipts/", handleReceiptRequests) // Unified handler for /receipts/{receipt_id} and /receipts/{receipt_id}/thumbnails

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

// handleReceiptRequests handles both /receipts/{receipt_id} and /receipts/{receipt_id}/thumbnails
func handleReceiptRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the URL ends with "/thumbnails"
	if strings.HasSuffix(r.URL.Path, "/thumbnails") {
		// Handle the thumbnail request
		handlers.GetThumbnails(w, r)
		return
	}

	// Otherwise, handle the receipt retrieval
	handlers.GetReceipt(w, r)
}
