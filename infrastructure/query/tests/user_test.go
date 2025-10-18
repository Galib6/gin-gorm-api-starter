package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	queryiface "github.com/zetsux/gin-gorm-clean-starter/core/interface/query"
	repositoryiface "github.com/zetsux/gin-gorm-clean-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/query"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	support "github.com/zetsux/gin-gorm-clean-starter/tests/testutil"
	"github.com/zetsux/gin-gorm-clean-starter/tests/testutil/factory"
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
	users, _, total, err := uq.GetAllUsers(context.Background(), base.GetsRequest{})
	require.NoError(t, err)
	require.GreaterOrEqual(t, int(total), 15)
	require.GreaterOrEqual(t, len(users), 15)

	// with pagination
	req := base.GetsRequest{Page: 1, PerPage: 10}
	users, last, total, err := uq.GetAllUsers(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, int64(2), last)
	require.Equal(t, int64(15), total)
	require.Equal(t, 10, len(users))
}
