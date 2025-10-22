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

func TestUserRepository_CreateUser(t *testing.T) {
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
	require.Equal(t, user.ID, created.ID)
	require.Equal(t, user.Email, created.Email)
	require.Equal(t, user.Name, created.Name)
	require.Equal(t, user.Role, created.Role)
}

func TestUserRepository_GetUserByPrimaryKey(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUser(t, ur, "New User", "newuser@mail.com", "passwordnew", constant.EnumRoleUser)

	fetched, err := ur.GetUserByPrimaryKey(ctx, nil, constant.DBAttrEmail, seed.Email)
	require.NoError(t, err)
	require.Equal(t, seed.ID, fetched.ID)
	require.Equal(t, seed.Email, fetched.Email)
	require.Equal(t, seed.Name, fetched.Name)
	require.Equal(t, seed.Role, fetched.Role)
}

func TestUserRepository_UpdateUser(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	newEmail := "newmail@example.com"
	edit := entity.User{ID: seed.ID, Email: newEmail}
	err := ur.UpdateUser(ctx, nil, edit)
	require.NoError(t, err)

	updated, err := ur.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, seed.ID.String())
	require.NoError(t, err)
	require.Equal(t, newEmail, updated.Email)
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	ur, ctx := setupUserRepositoryTest(t)

	seed := factory.SeedUsers(t, ur, 1)[0]

	err := ur.DeleteUserByID(ctx, nil, seed.ID.String())
	require.NoError(t, err)

	_, err = ur.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, seed.ID.String())
	require.Error(t, err)
}
