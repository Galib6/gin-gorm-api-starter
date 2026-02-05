package errs

import "errors"

var (
	// Product errors
	ErrProductNotFound      = errors.New("product not found")
	ErrProductSKUExists     = errors.New("product SKU already exists")
	ErrProductNoImage       = errors.New("product doesn't have any image")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidPriceRange    = errors.New("invalid price range")
	ErrInvalidStockQuantity = errors.New("invalid stock quantity")

	// Category errors
	ErrCategoryNotFound    = errors.New("category not found")
	ErrCategoryNameExists  = errors.New("category name already exists")
	ErrCategoryHasProducts = errors.New("category has associated products")
)
