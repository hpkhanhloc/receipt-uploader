package handlers

import (
	"fmt"
	"net/http"
	"receipt-uploader/models"
	"receipt-uploader/services"
	"strings"
)

// UploadReceipt handles the uploading of receipt images
func UploadReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the request's content type is multipart/form-data
	if r.Header.Get("Content-Type") == "" || !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	// Simulate user ID
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "test-user"
	}

	// Save the file using the service layer
	filePath, err := services.SaveFile(r)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}

	// Generate a unique receipt ID and store the receipt metadata
	receiptID := services.GenerateReceiptID()
	models.StoreReceipt(receiptID, filePath, userID)

	// Return the receipt ID
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Receipt uploaded successfully with ID: %s", receiptID)))
}
