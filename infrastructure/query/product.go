package query

import (
	"context"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	"myapp/support/base"

	"gorm.io/gorm"
)

var productAllowedSorts = []string{"id", "name", "sku", "price", "stock", "created_at", "updated_at"}
var productAllowedIncludes = []string{"Category"}

type productQuery struct {
	db *gorm.DB
}

func NewProductQuery(db *gorm.DB) *productQuery {
	return &productQuery{db: db}
}

// GetAllProducts returns products with complex filtering, sorting, and pagination
func (qr *productQuery) GetAllProducts(ctx context.Context, req dto.ProductGetsRequest,
) ([]entity.Product, base.PaginationResponse, error) {
	stmt := qr.db.WithContext(ctx).Debug().Model(&entity.Product{})

	// Filter by ID
	if req.ID != "" {
		stmt = stmt.Where("id = ?", req.ID)
	}

	// Filter by category
	if req.CategoryID != "" {
		stmt = stmt.Where("category_id = ?", req.CategoryID)
	}

	// Filter by active status
	if req.IsActive != nil {
		stmt = stmt.Where("is_active = ?", *req.IsActive)
	}

	// Search by name, description, or SKU
	if req.Search != "" {
		search := "%" + req.Search + "%"
		stmt = stmt.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", search, search, search)
	}

	// Price range filters
	if req.MinPrice != nil {
		stmt = stmt.Where("price >= ?", *req.MinPrice)
	}
	if req.MaxPrice != nil {
		stmt = stmt.Where("price <= ?", *req.MaxPrice)
	}

	// Stock range filters
	if req.MinStock != nil {
		stmt = stmt.Where("stock >= ?", *req.MinStock)
	}
	if req.MaxStock != nil {
		stmt = stmt.Where("stock <= ?", *req.MaxStock)
	}

	products, pageResp, err := GetWithPagination[entity.Product](stmt,
		req.PaginationRequest, productAllowedSorts, productAllowedIncludes)
	if err != nil {
		return nil, pageResp, err
	}
	return products, pageResp, nil
}

// GetProductsByPriceRange returns all products within a specific price range
func (qr *productQuery) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]entity.Product, error) {
	var products []entity.Product

	err := qr.db.WithContext(ctx).Debug().
		Model(&entity.Product{}).
		Where("price BETWEEN ? AND ?", minPrice, maxPrice).
		Where("is_active = ?", true).
		Order("price ASC").
		Find(&products).Error

	return products, err
}

// GetLowStockProducts returns products with stock below the threshold
func (qr *productQuery) GetLowStockProducts(ctx context.Context, threshold int) ([]entity.Product, error) {
	var products []entity.Product

	err := qr.db.WithContext(ctx).Debug().
		Model(&entity.Product{}).
		Where("stock < ?", threshold).
		Where("is_active = ?", true).
		Preload("Category").
		Order("stock ASC").
		Find(&products).Error

	return products, err
}

// GetProductStatsByCategory returns aggregated statistics per category
func (qr *productQuery) GetProductStatsByCategory(ctx context.Context) ([]dto.CategoryProductStats, error) {
	var stats []dto.CategoryProductStats

	err := qr.db.WithContext(ctx).Debug().
		Model(&entity.Product{}).
		Select(`
			categories.id AS category_id,
			categories.name AS category_name,
			COUNT(products.id) AS product_count,
			COALESCE(SUM(products.stock), 0) AS total_stock,
			COALESCE(AVG(products.price), 0) AS avg_price,
			COALESCE(MIN(products.price), 0) AS min_price,
			COALESCE(MAX(products.price), 0) AS max_price
		`).
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Where("products.deleted_at IS NULL").
		Group("categories.id, categories.name").
		Scan(&stats).Error

	return stats, err
}
