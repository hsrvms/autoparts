// internal/modules/inventory/services/service.go
package services

import (
	"context"
	"errors"

	inventorymodels "github.com/hsrvms/autoparts/internal/modules/inventory/models"
	"github.com/hsrvms/autoparts/internal/modules/inventory/repositories"
)

var (
	ErrItemNotFound        = errors.New("item not found")
	ErrDuplicatePartNumber = errors.New("part number already exists")
	ErrDuplicateBarcode    = errors.New("barcode already exists")
	ErrInvalidItemID       = errors.New("invalid item ID")
	ErrInvalidSubmodelID   = errors.New("invalid submodel ID")
	ErrCompatibilityExists = errors.New("compatibility already exists")
	ErrInvalidPrice        = errors.New("price must be greater than 0")
	ErrInvalidStock        = errors.New("stock cannot be negative")
)

type InventoryService interface {
	// Item operations
	GetItems(ctx context.Context, filter *inventorymodels.ItemFilter) ([]*inventorymodels.Item, error)
	GetItemByID(ctx context.Context, id int) (*inventorymodels.Item, error)
	GetItemByPartNumber(ctx context.Context, partNumber string) (*inventorymodels.Item, error)
	GetItemByBarcode(ctx context.Context, barcode string) (*inventorymodels.Item, error)
	CreateItem(ctx context.Context, item *inventorymodels.Item) (int, error)
	UpdateItem(ctx context.Context, item *inventorymodels.Item) error
	DeleteItem(ctx context.Context, id int) error
	GetLowStockItems(ctx context.Context) ([]*inventorymodels.Item, error)

	// Compatibility operations
	GetCompatibilities(ctx context.Context, itemID int) ([]*inventorymodels.Compatibility, error)
	AddCompatibility(ctx context.Context, compatibility *inventorymodels.Compatibility) (int, error)
	RemoveCompatibility(ctx context.Context, itemID, submodelID int) error
	GetCompatibleItems(ctx context.Context, submodelID int) ([]*inventorymodels.Item, error)
}

type inventoryService struct {
	repo repositories.InventoryRepository
}

func NewInventoryService(repo repositories.InventoryRepository) InventoryService {
	return &inventoryService{
		repo: repo,
	}
}

// Item operations
func (s *inventoryService) GetItems(ctx context.Context, filter *inventorymodels.ItemFilter) ([]*inventorymodels.Item, error) {
	return s.repo.GetItems(ctx, filter)
}

func (s *inventoryService) GetItemByID(ctx context.Context, id int) (*inventorymodels.Item, error) {
	if id <= 0 {
		return nil, ErrInvalidItemID
	}

	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrItemNotFound
	}

	return item, nil
}

func (s *inventoryService) GetItemByPartNumber(ctx context.Context, partNumber string) (*inventorymodels.Item, error) {
	if partNumber == "" {
		return nil, errors.New("part number is required")
	}

	return s.repo.GetItemByPartNumber(ctx, partNumber)
}

func (s *inventoryService) GetItemByBarcode(ctx context.Context, barcode string) (*inventorymodels.Item, error) {
	if barcode == "" {
		return nil, errors.New("barcode is required")
	}

	return s.repo.GetItemByBarcode(ctx, barcode)
}

func (s *inventoryService) CreateItem(ctx context.Context, item *inventorymodels.Item) (int, error) {
	// Validate required fields
	if err := s.validateItem(item); err != nil {
		return 0, err
	}

	// Check for duplicate part number
	existing, err := s.repo.GetItemByPartNumber(ctx, item.PartNumber)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, ErrDuplicatePartNumber
	}

	// Check for duplicate barcode if provided
	if item.Barcode != nil && *item.Barcode != "" {
		existing, err = s.repo.GetItemByBarcode(ctx, *item.Barcode)
		if err != nil {
			return 0, err
		}
		if existing != nil {
			return 0, ErrDuplicateBarcode
		}
	}

	return s.repo.CreateItem(ctx, item)
}

func (s *inventoryService) UpdateItem(ctx context.Context, item *inventorymodels.Item) error {
	if item.ItemID <= 0 {
		return ErrInvalidItemID
	}

	// Validate required fields
	if err := s.validateItem(item); err != nil {
		return err
	}

	// Check if item exists
	existing, err := s.repo.GetItemByID(ctx, item.ItemID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrItemNotFound
	}

	// Check for duplicate part number if changed
	if item.PartNumber != existing.PartNumber {
		existingByPartNumber, err := s.repo.GetItemByPartNumber(ctx, item.PartNumber)
		if err != nil {
			return err
		}
		if existingByPartNumber != nil && existingByPartNumber.ItemID != item.ItemID {
			return ErrDuplicatePartNumber
		}
	}

	// Check for duplicate barcode if changed
	if item.Barcode != nil && *item.Barcode != "" &&
		(existing.Barcode == nil || *item.Barcode != *existing.Barcode) {
		existingByBarcode, err := s.repo.GetItemByBarcode(ctx, *item.Barcode)
		if err != nil {
			return err
		}
		if existingByBarcode != nil && existingByBarcode.ItemID != item.ItemID {
			return ErrDuplicateBarcode
		}
	}

	return s.repo.UpdateItem(ctx, item)
}

func (s *inventoryService) DeleteItem(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidItemID
	}

	// Check if item exists
	existing, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrItemNotFound
	}

	// You might want to add additional checks here:
	// - Check if item has any active sales/purchases
	// - Check if item is part of any active orders
	// - Instead of deleting, maybe just mark as inactive

	return s.repo.DeleteItem(ctx, id)
}

func (s *inventoryService) GetLowStockItems(ctx context.Context) ([]*inventorymodels.Item, error) {
	return s.repo.GetLowStockItems(ctx)
}

// Compatibility operations
func (s *inventoryService) GetCompatibilities(ctx context.Context, itemID int) ([]*inventorymodels.Compatibility, error) {
	if itemID <= 0 {
		return nil, ErrInvalidItemID
	}

	return s.repo.GetCompatibilities(ctx, itemID)
}

func (s *inventoryService) AddCompatibility(ctx context.Context, compatibility *inventorymodels.Compatibility) (int, error) {
	if compatibility.ItemID <= 0 {
		return 0, ErrInvalidItemID
	}
	if compatibility.SubmodelID <= 0 {
		return 0, ErrInvalidSubmodelID
	}

	// Check if item exists
	item, err := s.repo.GetItemByID(ctx, compatibility.ItemID)
	if err != nil {
		return 0, err
	}
	if item == nil {
		return 0, ErrItemNotFound
	}

	// Check if compatibility already exists
	compatibilities, err := s.repo.GetCompatibilities(ctx, compatibility.ItemID)
	if err != nil {
		return 0, err
	}

	for _, existing := range compatibilities {
		if existing.SubmodelID == compatibility.SubmodelID {
			return 0, ErrCompatibilityExists
		}
	}

	return s.repo.AddCompatibility(ctx, compatibility)
}

func (s *inventoryService) RemoveCompatibility(ctx context.Context, itemID, submodelID int) error {
	if itemID <= 0 {
		return ErrInvalidItemID
	}
	if submodelID <= 0 {
		return ErrInvalidSubmodelID
	}

	return s.repo.RemoveCompatibility(ctx, itemID, submodelID)
}

func (s *inventoryService) GetCompatibleItems(ctx context.Context, submodelID int) ([]*inventorymodels.Item, error) {
	if submodelID <= 0 {
		return nil, ErrInvalidSubmodelID
	}

	return s.repo.GetCompatibleItems(ctx, submodelID)
}

// Helper functions
func (s *inventoryService) validateItem(item *inventorymodels.Item) error {
	if item.PartNumber == "" {
		return errors.New("part number is required")
	}
	if item.Description == "" {
		return errors.New("description is required")
	}
	if item.BuyPrice <= 0 {
		return errors.New("buy price must be greater than 0")
	}
	if item.SellPrice <= 0 {
		return errors.New("sell price must be greater than 0")
	}
	if item.CurrentStock < 0 {
		return errors.New("current stock cannot be negative")
	}
	if item.MinimumStock < 0 {
		return errors.New("minimum stock cannot be negative")
	}
	return nil
}
