package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/testutil"
	"github.com/zetsux/gin-gorm-clean-starter/testutil/factory"
)

func TestUserService_RegisterLoginFetchUpdateDelete(t *testing.T) {
	db := testutil.NewTestDB(t)
	us := service.NewUserService(factory.NewUserRepository(t, db))
	ctx := context.Background()

	// register
	created, err := us.CreateNewUser(ctx, dto.UserRegisterRequest{Name: "A", Email: "a@mail.test", Password: "secret"})
	require.NoError(t, err)
	require.Equal(t, "a@mail.test", created.Email)

	// login verify
	ok := us.VerifyLogin(ctx, "a@mail.test", "secret")
	require.True(t, ok)

	// fetch all
	users, page, err := us.GetAllUsers(ctx, base.GetsRequest{Page: 1, PerPage: 10})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 1)
	require.Equal(t, int64(1), page.LastPage)

	// get by email
	fetched, err := us.GetUserByPrimaryKey(ctx, "email", "a@mail.test")
	require.NoError(t, err)
	require.Equal(t, created.ID, fetched.ID)

	// update name
	edited, err := us.UpdateSelfName(ctx, dto.UserNameUpdateRequest{Name: "New"}, created.ID)
	require.NoError(t, err)
	require.Equal(t, "New", edited.Name)

	// delete
	err = us.DeleteUserByID(ctx, created.ID)
	require.NoError(t, err)
}
