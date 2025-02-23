package services

import (
	"context"
	"errors"

	suppliermodels "github.com/hsrvms/autoparts/internal/modules/suppliers/models"
	"github.com/hsrvms/autoparts/internal/modules/suppliers/repositories"
)

var (
	ErrSupplierNotFound      = errors.New("supplier not found")
	ErrInvalidSupplierID     = errors.New("invalid supplier ID")
	ErrDuplicateSupplierName = errors.New("supplier name already exists")
	ErrSupplierHasItems      = errors.New("cannot delete supplier with associated items")
)

type SupplierService interface {
	GetAll(ctx context.Context, filter *suppliermodels.SupplierFilter) ([]*suppliermodels.Supplier, error)
	GetByID(ctx context.Context, id int) (*suppliermodels.Supplier, error)
	Create(ctx context.Context, supplier *suppliermodels.Supplier) (int, error)
	Update(ctx context.Context, supplier *suppliermodels.Supplier) error
	Delete(ctx context.Context, id int) error
}

type supplierService struct {
	repo repositories.SupplierRepository
}

func NewSupplierService(repo repositories.SupplierRepository) SupplierService {
	return &supplierService{
		repo: repo,
	}
}

func (s *supplierService) GetAll(ctx context.Context, filter *suppliermodels.SupplierFilter) ([]*suppliermodels.Supplier, error) {
	return s.repo.GetAll(ctx, filter)
}

func (s *supplierService) GetByID(ctx context.Context, id int) (*suppliermodels.Supplier, error) {
	if id <= 0 {
		return nil, ErrInvalidSupplierID
	}

	supplier, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if supplier == nil {
		return nil, ErrSupplierNotFound
	}

	return supplier, nil
}

func (s *supplierService) Create(ctx context.Context, supplier *suppliermodels.Supplier) (int, error) {
	if err := s.validateSupplier(supplier); err != nil {
		return 0, err
	}

	// Check for existing suppliers with the same name
	existing, err := s.repo.GetAll(ctx, &suppliermodels.SupplierFilter{
		SearchTerm: &supplier.Name,
	})
	if err != nil {
		return 0, err
	}

	for _, s := range existing {
		if s.Name == supplier.Name {
			return 0, ErrDuplicateSupplierName
		}
	}

	return s.repo.Create(ctx, supplier)
}

func (s *supplierService) Update(ctx context.Context, supplier *suppliermodels.Supplier) error {
	if supplier.SupplierID <= 0 {
		return ErrInvalidSupplierID
	}

	if err := s.validateSupplier(supplier); err != nil {
		return err
	}

	// Check if supplier exists
	existing, err := s.repo.GetByID(ctx, supplier.SupplierID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSupplierNotFound
	}

	// Check for name uniqueness if name is being changed
	if existing.Name != supplier.Name {
		suppliers, err := s.repo.GetAll(ctx, &suppliermodels.SupplierFilter{
			SearchTerm: &supplier.Name,
		})
		if err != nil {
			return err
		}

		for _, s := range suppliers {
			if s.Name == supplier.Name && s.SupplierID != supplier.SupplierID {
				return ErrDuplicateSupplierName
			}
		}
	}

	return s.repo.Update(ctx, supplier)
}

func (s *supplierService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return ErrInvalidSupplierID
	}

	// Check if supplier exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSupplierNotFound
	}

	// Check if supplier has any associated items
	// Note: This would require a new repository method or a direct database check
	hasItems, err := s.checkSupplierHasItems(ctx, id)
	if err != nil {
		return err
	}
	if hasItems {
		return ErrSupplierHasItems
	}

	return s.repo.Delete(ctx, id)
}

// Helper functions
func (s *supplierService) validateSupplier(supplier *suppliermodels.Supplier) error {
	if supplier.Name == "" {
		return errors.New("supplier name is required")
	}
	// Add additional validations as needed
	return nil
}

// Note: This would need to be implemented with actual database access
func (s *supplierService) checkSupplierHasItems(ctx context.Context, supplierID int) (bool, error) {
	// For now, we'll return false to allow deletion
	// In a real implementation, you would check the items table
	return false, nil
}
