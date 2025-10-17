package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	repository_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-clean-starter/support/constant"
	support "github.com/zetsux/gin-gorm-clean-starter/tests/testutil"
	"github.com/zetsux/gin-gorm-clean-starter/tests/testutil/factory"
)

// --- Test Helpers ---

func setupUserRepositoryTest(t *testing.T) (repository_interface.UserRepository, context.Context) {
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

func TestUserRepository_UpdateName(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	newName := "New Name"
	updated, err := ur.UpdateNameUser(ctx, nil, newName, seed)
	require.NoError(t, err)
	require.Equal(t, newName, updated.Name)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	newEmail := "newmail@example.com"
	edit := entity.User{ID: seed.ID, Email: newEmail}
	edited, err := ur.UpdateUser(ctx, nil, edit)
	require.NoError(t, err)
	require.Equal(t, newEmail, edited.Email)
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	err := ur.DeleteUserByID(ctx, nil, seed.ID.String())
	require.NoError(t, err)
}
