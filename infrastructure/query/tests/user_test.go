package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	queryiface "github.com/zetsux/gin-gorm-api-starter/core/interface/query"
	repositoryiface "github.com/zetsux/gin-gorm-api-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-api-starter/infrastructure/query"
	"github.com/zetsux/gin-gorm-api-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	support "github.com/zetsux/gin-gorm-api-starter/tests/testutil"
	"github.com/zetsux/gin-gorm-api-starter/tests/testutil/factory"
)

// --- Test Helpers ---

func setupUserQueryTest(t *testing.T) (repositoryiface.UserRepository, queryiface.UserQuery, context.Context) {
	t.Helper()

	db := support.NewTestDB(t)
	ur := repository.NewUserRepository(db)
	uq := query.NewUserQuery(db)
	ctx := context.Background()

	return ur, uq, ctx
}

// --- Tests ---

func TestUserQuery_GetAllUsers_PaginationAndSearch(t *testing.T) {
	ur, uq, _ := setupUserQueryTest(t)

	_ = factory.SeedUsers(t, ur, 15)

	// no pagination
	users, pageResp, err := uq.GetAllUsers(context.Background(), dto.UserGetsRequest{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, int(pageResp.Total), 0)
	require.GreaterOrEqual(t, len(users), 15)

	// with pagination
	req := dto.UserGetsRequest{
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	}
	users, pageResp, err = uq.GetAllUsers(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, int64(2), pageResp.LastPage)
	require.Equal(t, int64(15), pageResp.Total)
	require.Equal(t, 10, len(users))
}
