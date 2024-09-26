package models

// Receipt represents the metadata of a receipt
type Receipt struct {
	ID       string
	FilePath string
	UserID   string
}

// In-memory receipt store
var receiptStore = make(map[string]Receipt)

// StoreReceipt saves the receipt metadata
func StoreReceipt(id, filePath, userID string) {
	receiptStore[id] = Receipt{
		ID:       id,
		FilePath: filePath,
		UserID:   userID,
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
