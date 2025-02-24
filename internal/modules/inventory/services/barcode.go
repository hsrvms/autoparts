package services

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"image/png"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// BarcodeService handles barcode operations
type BarcodeService interface {
	GenerateBarcode(categoryID int, supplierID int) (string, error)
	GenerateBarcodeImage(barcodeText string) ([]byte, error)
}

type barcodeService struct{}

func NewBarcodeService() BarcodeService {
	return &barcodeService{}
}

// GenerateBarcode creates a unique barcode for an item
func (s *barcodeService) GenerateBarcode(categoryID int, supplierID int) (string, error) {
	timestamp := time.Now().Format("060102150405")

	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	random := base32.StdEncoding.EncodeToString(randomBytes)[:8]

	barcode := fmt.Sprintf("C%03d-S%03d-%s-%s",
		categoryID,
		supplierID,
		timestamp,
		random,
	)

	return barcode, nil
}

// GenerateBarcodeImage generates a PNG image of the barcode
func (s *barcodeService) GenerateBarcodeImage(barcodeText string) ([]byte, error) {
	// Generate the barcode
	bc, err := code128.Encode(barcodeText)
	if err != nil {
		return nil, err
	}

	// Scale the barcode to a reasonable size
	scaled, err := barcode.Scale(bc, 300, 100)
	if err != nil {
		return nil, err
	}

	// Encode to PNG
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, scaled); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
