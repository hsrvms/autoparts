package repositories

import (
	"context"

	inventorymodels "github.com/hsrvms/autoparts/internal/modules/inventory/models"
)

type InventoryRepository interface {
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
