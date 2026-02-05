package dto

import (
	"mime/multipart"

	"myapp/support/base"

	"github.com/shopspring/decimal"
)

// ============== Product DTOs ==============

type (
	ProductGetsRequest struct {
		ID         string `json:"filter[id]" form:"filter[id]"`
		CategoryID string `json:"filter[category_id]" form:"filter[category_id]"`
		IsActive   *bool  `json:"filter[is_active]" form:"filter[is_active]"`
		Search     string `json:"search" form:"search"`

		// Price range filters
		MinPrice *float64 `json:"filter[min_price]" form:"filter[min_price]"`
		MaxPrice *float64 `json:"filter[max_price]" form:"filter[max_price]"`

		// Stock filters
		MinStock *int `json:"filter[min_stock]" form:"filter[min_stock]"`
		MaxStock *int `json:"filter[max_stock]" form:"filter[max_stock]"`

		base.PaginationRequest
	}

	ProductCreateRequest struct {
		Name        string  `json:"name" form:"name" binding:"required"`
		Description string  `json:"description" form:"description"`
		SKU         string  `json:"sku" form:"sku" binding:"required"`
		Price       float64 `json:"price" form:"price" binding:"required,gt=0"`
		Stock       int     `json:"stock" form:"stock" binding:"omitempty,min=0"`
		CategoryID  string  `json:"category_id" form:"category_id"`
		IsActive    *bool   `json:"is_active" form:"is_active"`
	}

	ProductUpdateRequest struct {
		ID          string   `json:"id"`
		Name        string   `json:"name" form:"name"`
		Description string   `json:"description" form:"description"`
		SKU         string   `json:"sku" form:"sku"`
		Price       *float64 `json:"price" form:"price" binding:"omitempty,gt=0"`
		Stock       *int     `json:"stock" form:"stock" binding:"omitempty,min=0"`
		CategoryID  string   `json:"category_id" form:"category_id"`
		IsActive    *bool    `json:"is_active" form:"is_active"`
	}

	ProductChangeImageRequest struct {
		ID    string                `json:"id"`
		Image *multipart.FileHeader `json:"image" form:"image"`
	}

	ProductStockUpdateRequest struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity" form:"quantity" binding:"required"`
	}

	ProductResponse struct {
		ID          string            `json:"id"`
		Name        string            `json:"name,omitempty"`
		Description string            `json:"description,omitempty"`
		SKU         string            `json:"sku,omitempty"`
		Price       decimal.Decimal   `json:"price,omitempty"`
		Stock       int               `json:"stock,omitempty"`
		CategoryID  string            `json:"category_id,omitempty"`
		IsActive    bool              `json:"is_active"`
		Image       string            `json:"image,omitempty"`
		Category    *CategoryResponse `json:"category,omitempty"`
	}
)

// ============== Category DTOs ==============

type (
	CategoryGetsRequest struct {
		ID     string `json:"filter[id]" form:"filter[id]"`
		Search string `json:"search" form:"search"`
		base.PaginationRequest
	}

	CategoryCreateRequest struct {
		Name        string `json:"name" form:"name" binding:"required"`
		Description string `json:"description" form:"description"`
	}

	CategoryUpdateRequest struct {
		ID          string `json:"id"`
		Name        string `json:"name" form:"name"`
		Description string `json:"description" form:"description"`
	}

	CategoryResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
	}

	CategoryProductStats struct {
		CategoryID   string  `json:"category_id"`
		CategoryName string  `json:"category_name"`
		ProductCount int64   `json:"product_count"`
		TotalStock   int64   `json:"total_stock"`
		AvgPrice     float64 `json:"avg_price"`
		MinPrice     float64 `json:"min_price"`
		MaxPrice     float64 `json:"max_price"`
	}
)

// ============== Product Maintenance DTOs ==============

// ProductMaintenanceRequest represents a complex, batch operation over products
type ProductMaintenanceRequest struct {
	ProductGetsRequest

	// Business rules / operations:

	// PriceAdjustment: percentage to adjust prices (e.g., 10 = +10%, -10 = -10%)
	PriceAdjustment *float64 `json:"price_adjustment" form:"price_adjustment"`

	// SetActive: if not nil, set all matching products to this active status
	SetActive *bool `json:"set_active" form:"set_active"`

	// LowStockThreshold: mark products with stock below this as inactive
	LowStockThreshold *int `json:"low_stock_threshold" form:"low_stock_threshold"`

	// NewCategoryID: move all matching products to this category
	NewCategoryID string `json:"new_category_id" form:"new_category_id"`
}

type ProductMaintenanceResult struct {
	ProductID       string   `json:"product_id"`
	ProductName     string   `json:"product_name"`
	OldPrice        *float64 `json:"old_price,omitempty"`
	NewPrice        *float64 `json:"new_price,omitempty"`
	ActiveChanged   bool     `json:"active_changed"`
	CategoryChanged bool     `json:"category_changed"`
}

type ProductMaintenanceResponse struct {
	TotalSelected        int                        `json:"total_selected"`
	TotalProcessed       int                        `json:"total_processed"`
	PriceChangedCount    int                        `json:"price_changed_count"`
	ActiveChangedCount   int                        `json:"active_changed_count"`
	CategoryChangedCount int                        `json:"category_changed_count"`
	Details              []ProductMaintenanceResult `json:"details"`
	base.PaginationResponse
}
