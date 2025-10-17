package factory

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	query_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/query"
	repository_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/query"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-clean-starter/support/constant"

	"gorm.io/gorm"
)

func NewUserRepository(t *testing.T, db *gorm.DB) repository_interface.UserRepository {
	t.Helper()
	return repository.NewUserRepository(db)
}

func NewUserQuery(t *testing.T, db *gorm.DB) query_interface.UserQuery {
	t.Helper()
	return query.NewUserQuery(db)
}

func NewUserService(t *testing.T, db *gorm.DB) service.UserService {
	t.Helper()
	ur := NewUserRepository(t, db)
	uq := NewUserQuery(t, db)
	return service.NewUserService(ur, uq)
}

func SeedUsers(t *testing.T, ur repository_interface.UserRepository, n int) []entity.User {
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

func SeedUser(t *testing.T, ur repository_interface.UserRepository, name, email, password, role string) entity.User {
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
