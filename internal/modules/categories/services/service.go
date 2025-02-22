package services

import (
	"context"

	"github.com/hsrvms/autoparts/internal/modules/categories/models"
	"github.com/hsrvms/autoparts/internal/modules/categories/repositories"
)

// CategoryService defines the interface for category business operations
type CategoryService interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*models.Category, error)
	GetSubcategories(ctx context.Context, parentID int) ([]*models.Category, error)
	CreateCategory(ctx context.Context, category *models.Category) (int, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id int) error
	GetCategoryTree(ctx context.Context) ([]*models.CategoryTreeNode, error)
}

// categoryService implements CategoryService
type categoryService struct {
	repo repositories.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

// GetAllCategories returns all categories
func (s *categoryService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	return s.repo.GetAll(ctx)
}

// GetCategoryByID returns a category by its ID
func (s *categoryService) GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// If category exists, get its subcategories
	if category != nil {
		subcategories, err := s.repo.GetSubcategories(ctx, category.CategoryID)
		if err != nil {
			return nil, err
		}
		category.Subcategories = subcategories
	}

	return category, nil
}

// GetSubcategories returns subcategories for a parent category
func (s *categoryService) GetSubcategories(ctx context.Context, parentID int) ([]*models.Category, error) {
	return s.repo.GetSubcategories(ctx, parentID)
}

// CreateCategory creates a new category
func (s *categoryService) CreateCategory(ctx context.Context, category *models.Category) (int, error) {
	// Add any business logic/validation here before calling repository

	// For example, we might want to check if the parent category exists
	if category.ParentCategoryID != nil {
		parent, err := s.repo.GetByID(ctx, *category.ParentCategoryID)
		if err != nil {
			return 0, err
		}
		if parent == nil {
			return 0, ErrParentCategoryNotFound
		}
	}

	return s.repo.Create(ctx, category)
}

// UpdateCategory updates an existing category
func (s *categoryService) UpdateCategory(ctx context.Context, category *models.Category) error {
	// Check if category exists
	existingCategory, err := s.repo.GetByID(ctx, category.CategoryID)
	if err != nil {
		return err
	}
	if existingCategory == nil {
		return ErrCategoryNotFound
	}

	// Check if parent category exists (if one is specified)
	if category.ParentCategoryID != nil {
		parent, err := s.repo.GetByID(ctx, *category.ParentCategoryID)
		if err != nil {
			return err
		}
		if parent == nil {
			return ErrParentCategoryNotFound
		}

		// Prevent circular references
		if *category.ParentCategoryID == category.CategoryID {
			return ErrCircularReference
		}
	}

	return s.repo.Update(ctx, category)
}

// DeleteCategory removes a category
func (s *categoryService) DeleteCategory(ctx context.Context, id int) error {
	// Check if category exists
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return ErrCategoryNotFound
	}

	// Check if category has subcategories
	subcategories, err := s.repo.GetSubcategories(ctx, id)
	if err != nil {
		return err
	}
	if len(subcategories) > 0 {
		return ErrCategoryHasSubcategories
	}

	return s.repo.Delete(ctx, id)
}

// GetCategoryTree returns a hierarchical structure of categories
func (s *categoryService) GetCategoryTree(ctx context.Context) ([]*models.CategoryTreeNode, error) {
	return s.repo.GetCategoryTree(ctx)
}
