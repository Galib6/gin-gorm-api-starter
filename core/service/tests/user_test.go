package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/core/entity"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-api-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-api-starter/core/service"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	"github.com/zetsux/gin-gorm-api-starter/support/util"
	"gorm.io/gorm"
)

// --- Mock Repository ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) DB() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockUserRepository) CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.User, error) {
	args := m.Called(ctx, tx, key, val)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) error {
	args := m.Called(ctx, tx, user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

// --- Mock Query ---

type MockUserQuery struct {
	mock.Mock
}

func (m *MockUserQuery) GetAllUsers(ctx context.Context, req dto.UserGetsRequest) ([]entity.User, base.PaginationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]entity.User), args.Get(1).(base.PaginationResponse), args.Error(2)
}

// --- Mock TxRepository (if needed) ---

type MockTxRepository struct {
	mock.Mock
}

func (m *MockTxRepository) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockTxRepository) Commit(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

func (m *MockTxRepository) Rollback(tx *gorm.DB) error {
	args := m.Called(tx)
	return args.Error(0)
}

// --- Test Helpers ---

func setupUserServiceMock() (service.UserService, *MockUserRepository, *MockUserQuery, context.Context) {
	repo := new(MockUserRepository)
	query := new(MockUserQuery)
	us := service.NewUserService(repo, query)
	ctx := context.Background()

	return us, repo, query, ctx
}

// --- Tests ---

func TestUserService_CreateNewUser(t *testing.T) {
	us, repo, _, ctx := setupUserServiceMock()

	expected := entity.User{ID: uuid.New(), Name: "A", Email: "a@mail.test"}
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "email", "a@mail.test").Return(entity.User{}, errs.ErrUserNotFound)
	repo.On("CreateNewUser", ctx, (*gorm.DB)(nil), mock.AnythingOfType("entity.User")).Return(expected, nil)

	user, err := us.CreateNewUser(ctx, dto.UserRegisterRequest{
		Name:     "A",
		Email:    "a@mail.test",
		Password: "secret",
	})
	require.NoError(t, err)
	require.Equal(t, "a@mail.test", user.Email)
	repo.AssertExpectations(t)
}

func TestUserService_VerifyLogin(t *testing.T) {
	us, repo, _, ctx := setupUserServiceMock()

	hashed, _ := util.PasswordHash("secret")
	stored := entity.User{ID: uuid.New(), Email: "a@mail.test", Password: hashed}
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "email", "a@mail.test").Return(stored, nil)

	ok := us.VerifyLogin(ctx, "a@mail.test", "secret")
	require.True(t, ok)
	repo.AssertExpectations(t)
}

func TestUserService_GetAllUsers(t *testing.T) {
	us, _, query, ctx := setupUserServiceMock()

	users := []entity.User{{ID: uuid.New(), Name: "A", Email: "a@mail.test"}}
	query.On("GetAllUsers", ctx, mock.AnythingOfType("dto.UserGetsRequest")).Return(users, base.PaginationResponse{LastPage: 1, Total: 1}, nil)

	got, page, err := us.GetAllUsers(ctx, dto.UserGetsRequest{
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 10,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, int64(1), page.LastPage)
	query.AssertExpectations(t)
}

func TestUserService_GetUserByPrimaryKey(t *testing.T) {
	us, repo, _, ctx := setupUserServiceMock()

	userID := uuid.New()
	expected := entity.User{ID: userID, Name: "A", Email: "a@mail.test"}
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "email", "a@mail.test").Return(expected, nil)

	fetched, err := us.GetUserByPrimaryKey(ctx, "email", "a@mail.test")
	require.NoError(t, err)
	require.Equal(t, userID.String(), fetched.ID)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateSelfName(t *testing.T) {
	us, repo, _, ctx := setupUserServiceMock()

	userID := uuid.New()
	old := entity.User{ID: userID, Name: "A"}
	updated := entity.User{ID: userID, Name: "New"}

	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "id", userID.String()).Return(old, nil)
	repo.On("UpdateUser", ctx, (*gorm.DB)(nil), updated).Return(nil)

	edited, err := us.UpdateSelfName(ctx, dto.UserNameUpdateRequest{ID: userID.String(), Name: "New"})
	require.NoError(t, err)
	require.Equal(t, "New", edited.Name)
	repo.AssertExpectations(t)
}

func TestUserService_DeleteUserByID(t *testing.T) {
	us, repo, _, ctx := setupUserServiceMock()

	userID := uuid.New()
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "id", userID.String()).Return(entity.User{ID: userID}, nil)
	repo.On("DeleteUserByID", ctx, (*gorm.DB)(nil), userID.String()).Return(nil)

	err := us.DeleteUserByID(ctx, userID.String())
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
