package models

import (
	"encoding/json"
	"log"
	"os"
)

// Receipt represents the metadata of a receipt
type Receipt struct {
	ID       string
	FilePath string
	UserID   string
}

// In-memory receipt store
var receiptStore = make(map[string]Receipt)

// File where receipts are stored
const receiptFile = "receipts.json"

// SaveReceiptsToFile saves the current in-memory receiptStore to a JSON file
func SaveReceiptsToFile() error {
	data, err := json.MarshalIndent(receiptStore, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(receiptFile, data, 0644)
}

// LoadReceiptsFromFile loads the receipt data from a JSON file into memory (receiptStore)
func LoadReceiptsFromFile() error {
	if _, err := os.Stat(receiptFile); os.IsNotExist(err) {
		return nil // If the file doesn't exist, skip loading
	}
	data, err := os.ReadFile(receiptFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &receiptStore)
}

// StoreReceipt saves the receipt metadata
func StoreReceipt(id, filePath, userID string) {
	receiptStore[id] = Receipt{
		ID:       id,
		FilePath: filePath,
		UserID:   userID,
	}
	err := SaveReceiptsToFile()
	if err != nil {
		log.Println("Error saving receipts to file:", err)
	}
}

// GetReceipt retrieves a receipt by ID
func GetReceipt(id string) (Receipt, bool) {
	receipt, exists := receiptStore[id]
	return receipt, exists
}

// ListUserReceipts returns all receipts for a given user
func ListUserReceipts(userID string) []Receipt {
	var receipts []Receipt
	for _, receipt := range receiptStore {
		if receipt.UserID == userID {
			receipts = append(receipts, receipt)
		}
	}
	return receipts
}
