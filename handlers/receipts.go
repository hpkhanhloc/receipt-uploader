package handlers

import (
	"encoding/json"
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

	// Extract user ID from headers
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "X-User-ID header is required", http.StatusBadRequest)
		return
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

// GetReceipt retrieves a receipt by ID and serves the file if the user is authorized
func GetReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from headers
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "X-User-ID header is required", http.StatusBadRequest)
		return
	}

	// Extract the receipt ID from the URL path
	receiptID := strings.TrimPrefix(r.URL.Path, "/receipts/")
	receipt, exists := models.GetReceipt(receiptID)
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Check if the user owns the receipt
	if receipt.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Serve the receipt file
	http.ServeFile(w, r, receipt.FilePath)
}

// ListReceipts lists all receipts for the authenticated user
func ListReceipts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from the X-User-ID header
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "X-User-ID header is required", http.StatusBadRequest)
		return
	}

	// Get the list of receipts for the user
	receipts := models.ListUserReceipts(userID)
	if len(receipts) == 0 {
		http.Error(w, "No receipts found for this user", http.StatusNotFound)
		return
	}

	// Return the full list of receipts as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipts)
}
