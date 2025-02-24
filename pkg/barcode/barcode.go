package barcode

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Generator handles barcode generation
type Generator struct{}

// New creates a new barcode generator
func New() *Generator {
	return &Generator{}
}

// Generate creates a new barcode for an item
// Format: XXnnnnnnYYv where:
// XX: Category prefix (2 letters)
// nnnnnn: Item ID (6 digits, zero-padded)
// YY: Year (2 digits)
// v: Validation digit
func (g *Generator) Generate(categoryName string, itemID int) (string, error) {
	// Get category prefix (first 2 letters, uppercase)
	prefix := strings.ToUpper(categoryName)
	if len(prefix) < 2 {
		return "", fmt.Errorf("category name too short")
	}
	prefix = prefix[:2]

	// Format item ID to 6 digits
	itemNum := fmt.Sprintf("%06d", itemID)

	// Get year suffix (last 2 digits)
	year := fmt.Sprintf("%02d", time.Now().Year()%100)

	// Create base code
	baseCode := fmt.Sprintf("%s%s%s", prefix, itemNum, year)

	// Calculate check digit
	checkDigit := g.calculateCheckDigit(baseCode)

	// Return complete barcode
	return fmt.Sprintf("%s%d", baseCode, checkDigit), nil
}

// Validate checks if a barcode is valid
func (g *Generator) Validate(barcode string) bool {
	// Check format: 10 characters + 1 check digit
	if len(barcode) != 11 {
		return false
	}

	// Check format with regex
	match, _ := regexp.MatchString(`^[A-Z]{2}\d{6}\d{2}\d{1}$`, barcode)
	if !match {
		return false
	}

	// Validate check digit
	baseCode := barcode[:10]
	checkDigit := int(barcode[10] - '0')
	return g.calculateCheckDigit(baseCode) == checkDigit
}

// calculateCheckDigit calculates the check digit for a barcode
func (g *Generator) calculateCheckDigit(baseCode string) int {
	sum := 0
	for i, r := range baseCode {
		// Alternate between multiplying by 3 and 1
		if i%2 == 0 {
			sum += int(r) * 3
		} else {
			sum += int(r)
		}
	}
	return (10 - (sum % 10)) % 10
}
