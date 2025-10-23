package repository

import (
	"context"
	"errors"

	"github.com/zetsux/gin-gorm-api-starter/core/entity"
	errs "github.com/zetsux/gin-gorm-api-starter/core/helper/errors"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (rp *userRepository) DB() *gorm.DB {
	return rp.db
}

func (rp *userRepository) CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	if tx == nil {
		tx = rp.db
	}

	if err := tx.WithContext(ctx).Debug().Create(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (rp *userRepository) GetUserByPrimaryKey(ctx context.Context,
	tx *gorm.DB, key string, val string) (entity.User, error) {
	var user entity.User

	if tx == nil {
		tx = rp.db
	}

	err := tx.WithContext(ctx).Debug().Where(key+" = $1", val).Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errs.ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}

func (rp *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) error {
	if tx == nil {
		tx = rp.db
	}

	if err := tx.WithContext(ctx).Debug().Updates(&user).Error; err != nil {
		return err
	}
	return nil
}

func (rp *userRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error {
	if tx == nil {
		tx = rp.db
	}

	if err := tx.WithContext(ctx).Debug().Delete(&entity.User{}, &id).Error; err != nil {
		return err
	}
	return nil
}
