package services

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"
)

// generateBarcode creates a unique barcode for an item
func generateBarcode(categoryID int, supplierID int) (string, error) {
	// Create a timestamp component (YYMMss)
	timestamp := time.Now().Format("060405")

	// Generate 3 random bytes
	randomBytes := make([]byte, 3)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Encode random bytes to base32 and take first 4 characters
	random := base32.StdEncoding.EncodeToString(randomBytes)[:4]

	// Combine all parts: C(3)S(3)-TS(6)-RND(4)
	barcode := fmt.Sprintf("%03d%03d%s%s",
		categoryID,
		supplierID,
		timestamp,
		random,
	)

	return barcode, nil
}
