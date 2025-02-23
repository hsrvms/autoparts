package repositories

import (
	"context"

	salesmodels "github.com/hsrvms/autoparts/internal/modules/sales/models"
)

type SaleRepository interface {
    GetAll(ctx context.Context, filter *salesmodels.SaleFilter) ([]*salesmodels.Sale, error)
    GetByID(ctx context.Context, id int) (*salesmodels.Sale, error)
    Create(ctx context.Context, sale *salesmodels.Sale) (int, error)
    Update(ctx context.Context, sale *salesmodels.Sale) error
    Delete(ctx context.Context, id int) error
    GetByTransactionNumber(ctx context.Context, transactionNumber string) (*salesmodels.Sale, error)
    GetItemSales(ctx context.Context, itemID int) ([]*salesmodels.Sale, error)
    GetCustomerSales(ctx context.Context, customerEmail string) ([]*salesmodels.Sale, error)
}
