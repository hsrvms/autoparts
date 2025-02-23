package repositories

import (
	"context"

	purchasemodels "github.com/hsrvms/autoparts/internal/modules/purchases/models"
)

type PurchaseRepository interface {
	GetAll(ctx context.Context, filter *purchasemodels.PurchaseFilter) ([]*purchasemodels.Purchase, error)
	GetByID(ctx context.Context, id int) (*purchasemodels.Purchase, error)
	Create(ctx context.Context, purchase *purchasemodels.Purchase) (int, error)
	Update(ctx context.Context, purchase *purchasemodels.Purchase) error
	Delete(ctx context.Context, id int) error
	GetByInvoiceNumber(ctx context.Context, invoiceNumber string) (*purchasemodels.Purchase, error)
	GetSupplierPurchases(ctx context.Context, supplierID int) ([]*purchasemodels.Purchase, error)
	GetItemPurchases(ctx context.Context, itemID int) ([]*purchasemodels.Purchase, error)
}
