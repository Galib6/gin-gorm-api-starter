package repositoryiface

import (
	"context"

	"myapp/core/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	// db
	DB() *gorm.DB

	// Product CRUD
	CreateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) (entity.Product, error)
	GetProductByID(ctx context.Context, tx *gorm.DB, id string, includes ...string) (entity.Product, error)
	GetProductByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.Product, error)
	UpdateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) error
	DeleteProductByID(ctx context.Context, tx *gorm.DB, id string) error

	// Batch operations
	UpdateProductStock(ctx context.Context, tx *gorm.DB, id string, quantity int) error
	BulkUpdatePrices(ctx context.Context, tx *gorm.DB, ids []string, priceMultiplier float64) error
}

type CategoryRepository interface {
	// db
	DB() *gorm.DB

	// Category CRUD
	CreateCategory(ctx context.Context, tx *gorm.DB, category entity.Category) (entity.Category, error)
	GetCategoryByID(ctx context.Context, tx *gorm.DB, id string, includes ...string) (entity.Category, error)
	GetCategoryByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.Category, error)
	UpdateCategory(ctx context.Context, tx *gorm.DB, category entity.Category) error
	DeleteCategoryByID(ctx context.Context, tx *gorm.DB, id string) error
}
