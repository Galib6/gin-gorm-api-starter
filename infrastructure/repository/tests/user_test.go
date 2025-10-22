package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/core/entity"
	repositoryiface "github.com/zetsux/gin-gorm-api-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-api-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"
	support "github.com/zetsux/gin-gorm-api-starter/tests/testutil"
	"github.com/zetsux/gin-gorm-api-starter/tests/testutil/factory"
)

// --- Test Helpers ---

func setupUserRepositoryTest(t *testing.T) (repositoryiface.UserRepository, context.Context) {
	t.Helper()

	db := support.NewTestDB(t)
	ur := repository.NewUserRepository(db)
	ctx := context.Background()

	return ur, ctx
}

// --- Tests ---

func TestUserRepository_CreateAndGetByPK(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	user := entity.User{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret",
		Role:     constant.EnumRoleUser,
	}

	created, err := ur.CreateNewUser(ctx, nil, user)
	require.NoError(t, err)
	require.Equal(t, user.Email, created.Email)

	fetched, err := ur.GetUserByPrimaryKey(ctx, nil, constant.DBAttrEmail, user.Email)
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	newEmail := "newmail@example.com"
	edit := entity.User{ID: seed.ID, Email: newEmail}
	err := ur.UpdateUser(ctx, nil, edit)
	require.NoError(t, err)
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	err := ur.DeleteUserByID(ctx, nil, seed.ID.String())
	require.NoError(t, err)
}
