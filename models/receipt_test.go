package models

import (
	"os"
	"path/filepath"
	"testing"
)

// Setup function to create a temporary environment for the test
func setupTestEnv() (string, error) {
	tmpDir := filepath.Join(os.TempDir(), "receipt_tests")
	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	ReceiptFile = filepath.Join(tmpDir, "../test_receipts.json") // Use a test file for receipts
	return tmpDir, nil
}

// Cleanup function to remove the temporary test environment
func cleanupTestEnv(tmpDir string) {
	os.RemoveAll(tmpDir)
}

// TestSaveAndLoadReceipts tests the saving and loading of receipts to/from a file
func TestSaveAndLoadReceipts(t *testing.T) {
	// Setup the test environment
	tmpDir, err := setupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer cleanupTestEnv(tmpDir)

	// Add some receipts to the in-memory store and store to file
	StoreReceipt("1", "/path/to/receipt1.jpg", "user1")
	StoreReceipt("2", "/path/to/receipt2.jpg", "user2")

	// Retrieve the receipt and verify the details
	receipt, exists := GetReceipt("1")
	if !exists || receipt.UserID != "user1" || receipt.FilePath != "/path/to/receipt1.jpg" {
		t.Fatalf("Receipt was not stored or retrieved correctly")
	}

	// Clear the in-memory store and load receipts from file
	ReceiptStore = make(map[string]Receipt)
	if err := LoadReceiptsFromFile(); err != nil {
		t.Fatalf("Failed to load receipts from file: %v", err)
	}

	// Verify the receipts were loaded correctly
	receipt1, exists := GetReceipt("1")
	if !exists || receipt1.UserID != "user1" || receipt1.FilePath != "/path/to/receipt1.jpg" {
		t.Fatalf("Receipt 1 was not loaded correctly")
	}

	receipt2, exists := GetReceipt("2")
	if !exists || receipt2.UserID != "user2" || receipt2.FilePath != "/path/to/receipt2.jpg" {
		t.Fatalf("Receipt 2 was not loaded correctly")
	}
}

// TestListUserReceipts tests listing receipts for a specific user
func TestListUserReceipts(t *testing.T) {
	// Setup the test environment
	tmpDir, err := setupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer cleanupTestEnv(tmpDir)

	// Add some receipts
	StoreReceipt("1", "/path/to/receipt1.jpg", "user1")
	StoreReceipt("2", "/path/to/receipt2.jpg", "user1")
	StoreReceipt("3", "/path/to/receipt3.jpg", "user2")

	// List receipts for user1
	userReceipts := ListUserReceipts("user1")
	if len(userReceipts) != 2 {
		t.Fatalf("Expected 2 receipts for user1, got %d", len(userReceipts))
	}

	// Verify the details of the receipts
	if userReceipts[0].ID != "1" && userReceipts[1].ID != "2" {
		t.Fatalf("Incorrect receipts returned for user1")
	}
}

// TestLoadReceiptsFromFileMissing tests loading receipts from a missing file
func TestLoadReceiptsFromFileMissing(t *testing.T) {
	// Setup the test environment
	tmpDir, err := setupTestEnv()
	if err != nil {
		t.Fatalf("Failed to setup test environment: %v", err)
	}
	defer cleanupTestEnv(tmpDir)

	// Ensure the receipt file does not exist
	os.Remove(ReceiptFile)

	// Attempt to load receipts from a missing file (should not error)
	err = LoadReceiptsFromFile()
	if err != nil {
		t.Fatalf("Expected no error when loading from a missing file, got: %v", err)
	}
}
