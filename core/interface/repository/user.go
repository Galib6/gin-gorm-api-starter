package repository

import (
	"context"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	// db
	DB() *gorm.DB

	// functional
	CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	GetUserByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.User, error)
	UpdateNameUser(ctx context.Context, tx *gorm.DB, name string, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error)
	DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error
}
