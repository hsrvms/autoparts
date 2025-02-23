package services

import (
	"context"
	"errors"
	"time"

	purchasemodels "github.com/hsrvms/autoparts/internal/modules/purchases/models"
	"github.com/hsrvms/autoparts/internal/modules/purchases/repositories"
)

var (
	ErrPurchaseNotFound       = errors.New("purchase not found")
	ErrInvalidPurchaseID      = errors.New("invalid purchase ID")
	ErrInvalidSupplierID      = errors.New("invalid supplier ID")
	ErrInvalidItemID          = errors.New("invalid item ID")
	ErrInvalidQuantity        = errors.New("quantity must be greater than 0")
	ErrInvalidCostPerUnit     = errors.New("cost per unit must be greater than 0")
	ErrDuplicateInvoiceNumber = errors.New("invoice number already exists")
	ErrInvalidDate            = errors.New("purchase date cannot be in the future")
)

type PurchaseService interface {
	GetAll(ctx context.Context, filter *purchasemodels.PurchaseFilter) ([]*purchasemodels.Purchase, error)
	GetByID(ctx context.Context, id int) (*purchasemodels.Purchase, error)
	Create(ctx context.Context, purchase *purchasemodels.Purchase) (int, error)
	Update(ctx context.Context, purchase *purchasemodels.Purchase) error
	Delete(ctx context.Context, id int) error
	GetSupplierPurchases(ctx context.Context, supplierID int) ([]*purchasemodels.Purchase, error)
	GetItemPurchases(ctx context.Context, itemID int) ([]*purchasemodels.Purchase, error)
}

type purchaseService struct {
	repo repositories.PurchaseRepository
}

func NewPurchaseService(repo repositories.PurchaseRepository) PurchaseService {
	return &purchaseService{
		repo: repo,
	}
}

func (s *purchaseService) GetAll(ctx context.Context, filter *purchasemodels.PurchaseFilter) ([]*purchasemodels.Purchase, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s *purchaseService) GetByID(ctx context.Context, id int) (*purchasemodels.Purchase, error) {
	if id <= 0 {
		return nil, ErrInvalidPurchaseID
	}

	purchase, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if purchase == nil {
		return nil, ErrPurchaseNotFound
	}

	return purchase, nil
}

func (s *purchaseService) Create(ctx context.Context, purchase *purchasemodels.Purchase) (int, error) {
	// Validate the purchase
	if err := s.validatePurchase(purchase); err != nil {
		return 0, err
	}

	// Check if invoice number is unique if provided
	if purchase.InvoiceNumber != nil && *purchase.InvoiceNumber != "" {
		existing, err := s.repo.GetByInvoiceNumber(ctx, *purchase.InvoiceNumber)
		if err != nil {
			return 0, err
		}
		if existing != nil {
			return 0, ErrDuplicateInvoiceNumber
		}
	}

	// Set date to current time if not provided
	if purchase.Date.IsZero() {
		purchase.Date = time.Now()
	}

	// Calculate total cost if not provided
	if purchase.TotalCost == 0 {
		purchase.TotalCost = float64(purchase.Quantity) * purchase.CostPerUnit
	}

	return s.repo.Create(ctx, purchase)
}

func (s *purchaseService) Update(ctx context.Context, purchase *purchasemodels.Purchase) error {
	if purchase.PurchaseID <= 0 {
		return ErrInvalidPurchaseID
	}

	// Validate the purchase
	if err := s.validatePurchase(purchase); err != nil {
		return err
	}

	// Check if purchase exists
	existing, err := s.repo.GetByID(ctx, purchase.PurchaseID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPurchaseNotFound
	}

	// Check if invoice number is unique if changed
	if purchase.InvoiceNumber != nil && *purchase.InvoiceNumber != "" &&
		(existing.InvoiceNumber == nil || *purchase.InvoiceNumber != *existing.InvoiceNumber) {
		existingWithInvoice, err := s.repo.GetByInvoiceNumber(ctx, *purchase.InvoiceNumber)
		if err != nil {
			return err
		}
		if existingWithInvoice != nil && existingWithInvoice.PurchaseID != purchase.PurchaseID {
			return ErrDuplicateInvoiceNumber
		}
	}

	// Recalculate total cost
	purchase.TotalCost = float64(purchase.Quantity) * purchase.CostPerUnit

	return s.repo.Update(ctx, purchase)
}

func (s *purchaseService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidPurchaseID
	}

	// Check if purchase exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPurchaseNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *purchaseService) GetSupplierPurchases(ctx context.Context, supplierID int) ([]*purchasemodels.Purchase, error) {
	if supplierID <= 0 {
		return nil, ErrInvalidSupplierID
	}
	return s.repo.GetSupplierPurchases(ctx, supplierID)
}

func (s *purchaseService) GetItemPurchases(ctx context.Context, itemID int) ([]*purchasemodels.Purchase, error) {
	if itemID <= 0 {
		return nil, ErrInvalidItemID
	}
	return s.repo.GetItemPurchases(ctx, itemID)
}

// Helper functions
func (s *purchaseService) validatePurchase(purchase *purchasemodels.Purchase) error {
	if purchase.SupplierID <= 0 {
		return ErrInvalidSupplierID
	}
	if purchase.ItemID <= 0 {
		return ErrInvalidItemID
	}
	if purchase.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	if purchase.CostPerUnit <= 0 {
		return ErrInvalidCostPerUnit
	}
	if !purchase.Date.IsZero() && purchase.Date.After(time.Now()) {
		return ErrInvalidDate
	}
	return nil
}
