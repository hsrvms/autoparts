package repositories

import (
	"context"
	"errors"

	"github.com/hsrvms/autoparts/internal/modules/categories/models"
	"github.com/hsrvms/autoparts/pkg/db"
	"github.com/jackc/pgx/v5"
)

// PostgresCategoryRepository implements CategoryRepository for PostgreSQL
type PostgresCategoryRepository struct {
	db *db.Database
}

// NewPostgresCategoryRepository creates a new PostgreSQL repository
func NewPostgresCategoryRepository(database *db.Database) CategoryRepository {
	return &PostgresCategoryRepository{
		db: database,
	}
}

// GetAll retrieves all categories from the database
func (r *PostgresCategoryRepository) GetAll(ctx context.Context) ([]*models.Category, error) {
	query := `
		SELECT category_id, category_name, description, parent_category_id, created_at, updated_at
		FROM categories
		ORDER BY category_name
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []*models.Category{}
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryName,
			&category.Description,
			&category.ParentCategoryID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetByID retrieves a category by its ID
func (r *PostgresCategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	query := `
		SELECT category_id, category_name, description, parent_category_id, created_at, updated_at
		FROM categories
		WHERE category_id = $1
	`

	category := &models.Category{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&category.CategoryID,
		&category.CategoryName,
		&category.Description,
		&category.ParentCategoryID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // No category found
		}
		return nil, err
	}

	return category, nil
}

// GetSubcategories retrieves all subcategories for a parent category
func (r *PostgresCategoryRepository) GetSubcategories(ctx context.Context, parentID int) ([]*models.Category, error) {
	query := `
		SELECT category_id, category_name, description, parent_category_id, created_at, updated_at
		FROM categories
		WHERE parent_category_id = $1
		ORDER BY category_name
	`

	rows, err := r.db.Pool.Query(ctx, query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subcategories := []*models.Category{}
	for rows.Next() {
		category := &models.Category{}
		err := rows.Scan(
			&category.CategoryID,
			&category.CategoryName,
			&category.Description,
			&category.ParentCategoryID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subcategories = append(subcategories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subcategories, nil
}

// Create adds a new category to the database
func (r *PostgresCategoryRepository) Create(ctx context.Context, category *models.Category) (int, error) {
	query := `
		INSERT INTO categories (category_name, description, parent_category_id)
		VALUES ($1, $2, $3)
		RETURNING category_id
	`

	var id int
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		category.CategoryName,
		category.Description,
		category.ParentCategoryID,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update modifies an existing category
func (r *PostgresCategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories
		SET category_name = $2, description = $3, parent_category_id = $4
		WHERE category_id = $1
	`

	_, err := r.db.Pool.Exec(
		ctx,
		query,
		category.CategoryID,
		category.CategoryName,
		category.Description,
		category.ParentCategoryID,
	)

	return err
}

// Delete removes a category from the database
func (r *PostgresCategoryRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE category_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// GetCategoryTree builds a hierarchical tree of categories
func (r *PostgresCategoryRepository) GetCategoryTree(ctx context.Context) ([]*models.CategoryTreeNode, error) {
	// First, get all categories
	categories, err := r.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	categoryMap := make(map[int]*models.Category)
	for _, category := range categories {
		categoryMap[category.CategoryID] = category
	}

	// Build the tree
	rootNodes := []*models.CategoryTreeNode{}

	// Process each category
	for _, category := range categories {
		// If it's a root category (no parent)
		if category.ParentCategoryID == nil {
			node := &models.CategoryTreeNode{
				Category:      category,
				Subcategories: []*models.CategoryTreeNode{},
			}
			rootNodes = append(rootNodes, node)
		}
	}

	// Process child categories
	for _, rootNode := range rootNodes {
		buildSubtree(rootNode, categories)
	}

	return rootNodes, nil
}

// Helper function to build the category subtree
func buildSubtree(node *models.CategoryTreeNode, allCategories []*models.Category) {
	for _, category := range allCategories {
		if category.ParentCategoryID != nil && *category.ParentCategoryID == node.Category.CategoryID {
			childNode := &models.CategoryTreeNode{
				Category:      category,
				Subcategories: []*models.CategoryTreeNode{},
			}
			node.Subcategories = append(node.Subcategories, childNode)
			buildSubtree(childNode, allCategories)
		}
	}
}
