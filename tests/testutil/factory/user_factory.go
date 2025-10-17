package factory

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/support/constant"

	"gorm.io/gorm"
)

func NewUserRepository(t *testing.T, db *gorm.DB) repository.UserRepository {
	t.Helper()
	txr := repository.NewTxRepository(db)
	return repository.NewUserRepository(txr)
}

func NewUserService(t *testing.T, db *gorm.DB) service.UserService {
	t.Helper()
	ur := NewUserRepository(t, db)
	return service.NewUserService(ur)
}

func SeedUsers(t *testing.T, ur repository.UserRepository, n int) []entity.User {
	ctx := context.Background()
	users := make([]entity.User, 0, n)
	for i := 0; i < n; i++ {
		u := entity.User{
			ID:       uuid.New(),
			Name:     "User" + uuid.NewString()[:8],
			Email:    uuid.NewString() + "@mail.test",
			Password: "password",
			Role:     constant.EnumRoleUser,
		}
		created, err := ur.CreateNewUser(ctx, nil, u)
		require.NoError(t, err)
		users = append(users, created)
	}
	return users
}

func SeedUser(t *testing.T, ur repository.UserRepository, name, email, password, role string) entity.User {
	t.Helper()
	ctx := context.Background()
	u := entity.User{
		ID:       uuid.New(),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
	created, err := ur.CreateNewUser(ctx, nil, u)
	require.NoError(t, err)
	return created
}
