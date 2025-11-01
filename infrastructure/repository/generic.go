package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

func useDB(tx *gorm.DB, defaultDB *gorm.DB) *gorm.DB {
	if tx == nil {
		return defaultDB
	}
	return tx
}

func Create[T any](ctx context.Context, tx *gorm.DB, defaultDB *gorm.DB, entity T) (T, error) {
	if err := useDB(tx, defaultDB).WithContext(ctx).Debug().Create(&entity).Error; err != nil {
		return entity, err
	}
	return entity, nil
}

func GetByID[T any](ctx context.Context, tx *gorm.DB, defaultDB *gorm.DB,
	id string, notFoundErr error, includes ...string) (T, error) {
	var entity T

	stmt := useDB(tx, defaultDB).WithContext(ctx).Debug().Model(&entity)
	for _, include := range includes {
		stmt = stmt.Preload(include)
	}

	err := stmt.Where("id = ?", id).Take(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, notFoundErr
		}
		return entity, err
	}
	return entity, nil
}

func Update[T any](ctx context.Context, tx *gorm.DB, defaultDB *gorm.DB, entity *T) error {
	if err := useDB(tx, defaultDB).WithContext(ctx).Debug().Updates(entity).Error; err != nil {
		return err
	}
	return nil
}

func Delete[T any](ctx context.Context, tx *gorm.DB, defaultDB *gorm.DB, id string) error {
	obj := new(T)
	if err := useDB(tx, defaultDB).WithContext(ctx).Debug().Delete(obj, &id).Error; err != nil {
		return err
	}
	return nil
}
