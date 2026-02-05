package repository

import (
	"context"
	"errors"

	"myapp/core/entity"
	errs "myapp/core/helper/errors"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (rp *productRepository) DB() *gorm.DB {
	return rp.db
}

func (rp *productRepository) CreateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) (entity.Product, error) {
	return Create(ctx, tx, rp.DB(), product)
}

func (rp *productRepository) GetProductByID(ctx context.Context, tx *gorm.DB, id string, includes ...string) (entity.Product, error) {
	return GetByID[entity.Product](ctx, tx, rp.DB(), id, errs.ErrProductNotFound, includes...)
}

func (rp *productRepository) GetProductByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.Product, error) {
	var product entity.Product

	if tx == nil {
		tx = rp.db
	}

	err := tx.WithContext(ctx).Debug().Where(key+" = $1", val).Take(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Product{}, errs.ErrProductNotFound
		}
		return product, err
	}
	return product, nil
}

func (rp *productRepository) UpdateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) error {
	return Update(ctx, tx, rp.DB(), &product)
}

func (rp *productRepository) DeleteProductByID(ctx context.Context, tx *gorm.DB, id string) error {
	return Delete[entity.Product](ctx, tx, rp.DB(), id)
}

// UpdateProductStock updates the stock of a product by adding the given quantity
// Quantity can be negative for stock reduction
func (rp *productRepository) UpdateProductStock(ctx context.Context, tx *gorm.DB, id string, quantity int) error {
	db := useDB(tx, rp.db)

	result := db.WithContext(ctx).Debug().
		Model(&entity.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrProductNotFound
	}

	return nil
}

// BulkUpdatePrices updates the prices of multiple products using a multiplier
// e.g., multiplier of 1.10 increases prices by 10%
func (rp *productRepository) BulkUpdatePrices(ctx context.Context, tx *gorm.DB, ids []string, priceMultiplier float64) error {
	if len(ids) == 0 {
		return nil
	}

	db := useDB(tx, rp.db)

	result := db.WithContext(ctx).Debug().
		Model(&entity.Product{}).
		Where("id IN ?", ids).
		Update("price", gorm.Expr("price * ?", priceMultiplier))

	return result.Error
}
