package repository

import (
	"context"
	"errors"

	"myapp/core/entity"
	errs "myapp/core/helper/errors"

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

func (ur *userRepository) CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	return Create(ctx, tx, ur.DB(), user)
}

func (ur *userRepository) GetUserByPrimaryKey(ctx context.Context,
	tx *gorm.DB, key string, val string) (entity.User, error) {
	var user entity.User

	if tx == nil {
		tx = ur.db
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

func (ur *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) error {
	return Update(ctx, tx, ur.DB(), &user)
}

func (ur *userRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error {
	return Delete[entity.User](ctx, tx, ur.DB(), id)
}
