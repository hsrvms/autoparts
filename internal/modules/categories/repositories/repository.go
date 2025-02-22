package repositories

import (
	"context"

	"github.com/hsrvms/autoparts/internal/modules/categories/models"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]*models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	GetSubcategories(ctx context.Context, parentID int) ([]*models.Category, error)
	Create(ctx context.Context, category *models.Category) (int, error)
	Update(ctx context.Context, categoty *models.Category) error
	Delete(ctx context.Context, id int) error
	GetCategoryTree(ctx context.Context) ([]*models.CategoryTreeNode, error)
}
