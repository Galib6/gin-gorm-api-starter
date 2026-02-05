package repository

import (
	"context"
	"errors"

	"myapp/core/entity"
	errs "myapp/core/helper/errors"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db: db}
}

func (rp *categoryRepository) DB() *gorm.DB {
	return rp.db
}

func (rp *categoryRepository) CreateCategory(ctx context.Context, tx *gorm.DB, category entity.Category) (entity.Category, error) {
	return Create(ctx, tx, rp.DB(), category)
}

func (rp *categoryRepository) GetCategoryByID(ctx context.Context, tx *gorm.DB, id string, includes ...string) (entity.Category, error) {
	return GetByID[entity.Category](ctx, tx, rp.DB(), id, errs.ErrCategoryNotFound, includes...)
}

func (rp *categoryRepository) GetCategoryByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.Category, error) {
	var category entity.Category

	if tx == nil {
		tx = rp.db
	}

	err := tx.WithContext(ctx).Debug().Where(key+" = $1", val).Take(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Category{}, errs.ErrCategoryNotFound
		}
		return category, err
	}
	return category, nil
}

func (rp *categoryRepository) UpdateCategory(ctx context.Context, tx *gorm.DB, category entity.Category) error {
	return Update(ctx, tx, rp.DB(), &category)
}

func (rp *categoryRepository) DeleteCategoryByID(ctx context.Context, tx *gorm.DB, id string) error {
	return Delete[entity.Category](ctx, tx, rp.DB(), id)
}
