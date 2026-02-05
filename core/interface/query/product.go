package queryiface

import (
	"context"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	"myapp/support/base"
)

type ProductQuery interface {
	// Basic query with pagination
	GetAllProducts(ctx context.Context, req dto.ProductGetsRequest) ([]entity.Product, base.PaginationResponse, error)

	// Complex aggregated queries
	GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]entity.Product, error)
	GetLowStockProducts(ctx context.Context, threshold int) ([]entity.Product, error)
	GetProductStatsByCategory(ctx context.Context) ([]dto.CategoryProductStats, error)
}

type CategoryQuery interface {
	GetAllCategories(ctx context.Context, req dto.CategoryGetsRequest) ([]entity.Category, base.PaginationResponse, error)
}
