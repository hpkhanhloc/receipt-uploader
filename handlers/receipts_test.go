package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"receipt-uploader/models"
	"strings"
	"testing"
)

// Setup function to create a temporary environment for the test
func setupTestEnv(t *testing.T) string {
	tmpDir := t.TempDir()
	models.ReceiptFile = filepath.Join(tmpDir, "test_receipts.json")
	models.ReceiptStore = make(map[string]models.Receipt)
	return tmpDir
}

// TestUploadReceipt tests the UploadReceipt handler
func TestUploadReceipt(t *testing.T) {
	// Setup the in-memory store for the test and use temporary receipts.json
	setupTestEnv(t)
	models.StoreReceipt("1", "/path/to/receipt1.jpg", "test-user")

	t.Run("MissingUserIDHeader", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/receipts", strings.NewReader(""))
		req.Header.Set("Content-Type", "multipart/form-data")

		rr := httptest.NewRecorder()

		UploadReceipt(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", rr.Code)
		}
	})
}

// TestGetReceipt tests the GetReceipt handler
func TestGetReceipt(t *testing.T) {
	// Setup in-memory store with some sample receipts and use temporary receipts.json
	setupTestEnv(t)
	models.StoreReceipt("1", "../testdata/test.jpg", "test-user")

	t.Run("ValidGetReceipt", func(t *testing.T) {
		// Create a valid GET request
		req := httptest.NewRequest(http.MethodGet, "/receipts/1", nil)
		req.Header.Set("X-User-ID", "test-user")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		GetReceipt(rr, req)

		// Check the status code
		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", rr.Code)
		}
	})

	t.Run("ReceiptNotFound", func(t *testing.T) {
		// Create a GET request for a non-existent receipt
		req := httptest.NewRequest(http.MethodGet, "/receipts/999", nil)
		req.Header.Set("X-User-ID", "test-user")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		GetReceipt(rr, req)

		// Check the status code
		if rr.Code != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", rr.Code)
		}
	})

	t.Run("UnauthorizedAccess", func(t *testing.T) {
		// Create a GET request for an existing receipt with a different user
		req := httptest.NewRequest(http.MethodGet, "/receipts/1", nil)
		req.Header.Set("X-User-ID", "another-user")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		GetReceipt(rr, req)

		// Check the status code
		if rr.Code != http.StatusForbidden {
			t.Fatalf("Expected status code 403, got %d", rr.Code)
		}
	})

	t.Run("MissingUserIDHeader", func(t *testing.T) {
		// Create a GET request with no X-User-ID header
		req := httptest.NewRequest(http.MethodGet, "/receipts/1", nil)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		GetReceipt(rr, req)

		// Check the status code
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", rr.Code)
		}
	})

	t.Run("ValidImageResize", func(t *testing.T) {
		// Create a GET request with width and height query parameters
		req := httptest.NewRequest(http.MethodGet, "/receipts/1?width=100&height=100", nil)
		req.Header.Set("X-User-ID", "test-user")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		GetReceipt(rr, req)

		// Check the status code
		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", rr.Code)
		}

		// Optionally, check if the response Content-Type is image/jpeg
		if contentType := rr.Header().Get("Content-Type"); contentType != "image/jpeg" {
			t.Fatalf("Expected Content-Type image/jpeg, got %s", contentType)
		}
	})
}

// TestListReceipts tests the ListReceipts handler
func TestListReceipts(t *testing.T) {
	// Setup some sample receipts in the in-memory store and use temporary receipts.json
	setupTestEnv(t)
	models.StoreReceipt("1", "/path/to/receipt1.jpg", "test-user")
	models.StoreReceipt("2", "/path/to/receipt2.jpg", "test-user")
	models.StoreReceipt("3", "/path/to/receipt3.jpg", "another-user")

	t.Run("ValidListReceipts", func(t *testing.T) {
		// Create a GET request for listing receipts for the test-user
		req := httptest.NewRequest(http.MethodGet, "/receipts", nil)
		req.Header.Set("X-User-ID", "test-user")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		ListReceipts(rr, req)

		// Check the status code
		if rr.Code != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", rr.Code)
		}

		// Optionally, check the response body for expected JSON structure
		expected := `[{"ID":"1","FilePath":"/path/to/receipt1.jpg","UserID":"test-user"},{"ID":"2","FilePath":"/path/to/receipt2.jpg","UserID":"test-user"}]`
		if strings.TrimSpace(rr.Body.String()) != expected {
			t.Fatalf("Expected JSON: %s, got: %s", expected, rr.Body.String())
		}
	})

	t.Run("NoReceiptsForUser", func(t *testing.T) {
		// Create a GET request for a user with no receipts
		req := httptest.NewRequest(http.MethodGet, "/receipts", nil)
		req.Header.Set("X-User-ID", "empty-user")

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		ListReceipts(rr, req)

		// Check the status code
		if rr.Code != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", rr.Code)
		}
	})

	t.Run("MissingUserIDHeader", func(t *testing.T) {
		// Create a GET request with no X-User-ID header
		req := httptest.NewRequest(http.MethodGet, "/receipts", nil)

		// Create a response recorder
		rr := httptest.NewRecorder()

		// Call the handler
		ListReceipts(rr, req)

		// Check the status code
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", rr.Code)
		}
	})
}
