package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"receipt-uploader/models"
	"receipt-uploader/services"
	"strconv"
	"strings"
	"sync"

	"github.com/disintegration/imaging"
)

// ThumbnailResponse holds the paths for the generated thumbnails
type ThumbnailResponse struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

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

	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // Max 10MB
	if err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Retrieve all files from the form
	files := r.MultipartForm.File["file"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish
	receiptIDs := make([]string, len(files))
	errs := make([]error, len(files))

	// Process each file concurrently
	for i, fileHeader := range files {
		wg.Add(1)
		go func(i int, fileHeader *multipart.FileHeader) {
			defer wg.Done()

			// Save the file using the service layer
			filePath, err := services.SaveFile(fileHeader)
			if err != nil {
				errs[i] = err
				return
			}

			// Generate a unique receipt ID and store the receipt metadata
			receiptID := services.GenerateReceiptID()
			models.StoreReceipt(receiptID, filePath, userID)
			receiptIDs[i] = receiptID
		}(i, fileHeader)
	}

	// Wait for all the goroutines to finish
	wg.Wait()

	// Check if any errors occurred
	for _, err := range errs {
		if err != nil {
			http.Error(w, "Error uploading one or more files", http.StatusInternalServerError)
			return
		}
	}

	// Return the list of receipt IDs
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Receipts uploaded successfully with IDs: %v", strings.Join(receiptIDs, ", "))))
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

	// Parse optional width and height query parameters
	width, err := parseQueryParameter(r.URL.Query().Get("width"), "width")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	height, err := parseQueryParameter(r.URL.Query().Get("height"), "height")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If neither width nor height is provided, serve the original image
	if width == 0 && height == 0 {
		http.ServeFile(w, r, receipt.FilePath)
		return
	}

	// Process the image (resize based on provided width and height)
	img, err := services.ProcessImage(receipt.FilePath, width, height)
	if err != nil {
		http.Error(w, "Could not process image", http.StatusInternalServerError)
		return
	}

	// Serve the resized image back to the client
	w.Header().Set("Content-Type", "image/jpeg")
	err = imaging.Encode(w, img, imaging.JPEG)
	if err != nil {
		http.Error(w, "Could not encode resized image", http.StatusInternalServerError)
		return
	}
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

// parseQueryParameter parses a query parameter and returns its integer value
func parseQueryParameter(paramStr, paramName string) (int, error) {
	if paramStr == "" {
		return 0, nil // If the parameter is not provided, return 0
	}

	value, err := strconv.Atoi(paramStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s parameter", paramName)
	}

	if value < 0 {
		return 0, fmt.Errorf("%s must be a positive integer", paramName)
	}

	return value, nil
}

// GetThumbnails generates thumbnails in small, medium, and large sizes using ProcessImageConcurrently
func GetThumbnails(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from headers
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "X-User-ID header is required", http.StatusBadRequest)
		return
	}

	// Extract the receipt ID from the URL path
	receiptID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/receipts/"), "/thumbnails")
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

	// Use a WaitGroup to wait for concurrent processing
	var wg sync.WaitGroup
	wg.Add(3)

	// Channel to collect results
	resultCh := make(chan services.Result, 3)

	// Process the image concurrently for small, medium, and large sizes
	go func() {
		defer wg.Done()
		img, err := services.ProcessImage(receipt.FilePath, 100, 100) // Small
		resultCh <- services.Result{Img: img, Err: err}
	}()
	go func() {
		defer wg.Done()
		img, err := services.ProcessImage(receipt.FilePath, 200, 200) // Medium
		resultCh <- services.Result{Img: img, Err: err}
	}()
	go func() {
		defer wg.Done()
		img, err := services.ProcessImage(receipt.FilePath, 400, 400) // Large
		resultCh <- services.Result{Img: img, Err: err}
	}()

	// Wait for all processing to complete, then close resultCh
	go func() {
		wg.Wait()
		close(resultCh) // Only close the channel after all the Goroutines are done
	}()

	// Prepare the response
	thumbnailResponse := ThumbnailResponse{}

	// Collect results
	for res := range resultCh {
		if res.Err != nil {
			http.Error(w, fmt.Sprintf("Error processing image: %v", res.Err), http.StatusInternalServerError)
			return
		}

		width := res.Img.Bounds().Dx() // Get the width of the processed image
		// Save the resized image and determine which size it is based on its dimensions
		thumbnailPath := filepath.Join("uploads", fmt.Sprintf("%s_%dx%d.jpg", receiptID, res.Img.Bounds().Dx(), res.Img.Bounds().Dy()))

		// Save the image to the appropriate thumbnail size file
		if err := services.SaveImage(res.Img, thumbnailPath); err != nil {
			http.Error(w, fmt.Sprintf("Error saving thumbnail: %v", err), http.StatusInternalServerError)
			return
		}

		// Assign the path to the correct field in the response
		switch {
		case width <= 100:
			thumbnailResponse.Small = thumbnailPath
		case width <= 200:
			thumbnailResponse.Medium = thumbnailPath
		case width <= 400:
			thumbnailResponse.Large = thumbnailPath
		default:
			log.Printf("Unexpected image dimensions: %dx%d", res.Img.Bounds().Dx(), res.Img.Bounds().Dy())
		}

	}

	// Return the response with thumbnail paths
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(thumbnailResponse)
}
