package services

import "errors"

// Service errors
var (
	ErrCategoryNotFound         = errors.New("category not found")
	ErrParentCategoryNotFound   = errors.New("parent category not found")
	ErrCategoryHasSubcategories = errors.New("category has subcategories and cannot be deleted")
	ErrCircularReference        = errors.New("circular reference detected: a category cannot be its own parent")
)
