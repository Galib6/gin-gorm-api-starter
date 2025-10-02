package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/testutil"
)

func seedUsers(t *testing.T, ur repository.UserRepository, n int) []entity.User {
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

func TestUserRepository_CreateAndGetByPK(t *testing.T) {
	db := testutil.NewTestDB(t)
	ur := repository.NewUserRepository(repository.NewTxRepository(db))
	ctx := context.Background()

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

func TestUserRepository_UpdateNameAndUpdateUser(t *testing.T) {
	db := testutil.NewTestDB(t)
	ur := repository.NewUserRepository(repository.NewTxRepository(db))
	ctx := context.Background()

	seed := seedUsers(t, ur, 1)[0]

	newName := "New Name"
	updated, err := ur.UpdateNameUser(ctx, nil, newName, seed)
	require.NoError(t, err)
	require.Equal(t, newName, updated.Name)

	newEmail := "newmail@example.com"
	edit := entity.User{ID: updated.ID, Email: newEmail}
	edited, err := ur.UpdateUser(ctx, nil, edit)
	require.NoError(t, err)
	require.Equal(t, newEmail, edited.Email)
}

func TestUserRepository_GetAllUsers_PaginationAndSearch(t *testing.T) {
	db := testutil.NewTestDB(t)
	ur := repository.NewUserRepository(repository.NewTxRepository(db))
	ctx := context.Background()
	_ = ctx

	_ = seedUsers(t, ur, 15)

	// no pagination
	users, _, total, err := ur.GetAllUsers(context.Background(), nil, base.GetsRequest{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, int(total), 15)
	require.GreaterOrEqual(t, len(users), 15)

	// with pagination
	req := base.GetsRequest{Page: 1, PerPage: 10}
	users, last, total, err := ur.GetAllUsers(context.Background(), nil, req)
	require.NoError(t, err)
	require.Equal(t, int64(2), last)
	require.Equal(t, int64(15), total)
	require.Equal(t, 10, len(users))
}

func TestUserRepository_DeleteUserByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	ur := repository.NewUserRepository(repository.NewTxRepository(db))
	ctx := context.Background()

	seed := seedUsers(t, ur, 1)[0]

	err := ur.DeleteUserByID(ctx, nil, seed.ID.String())
	require.NoError(t, err)
}
