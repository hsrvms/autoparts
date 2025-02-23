package repositories

import (
	"context"

	suppliermodels "github.com/hsrvms/autoparts/internal/modules/suppliers/models"
)

type SupplierRepository interface {
    GetAll(ctx context.Context, filter *suppliermodels.SupplierFilter) ([]*suppliermodels.Supplier, error)
    GetByID(ctx context.Context, id int) (*suppliermodels.Supplier, error)
    Create(ctx context.Context, supplier *suppliermodels.Supplier) (int, error)
    Update(ctx context.Context, supplier *suppliermodels.Supplier) error
    Delete(ctx context.Context, id int) error
}
