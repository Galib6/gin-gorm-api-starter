package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	"github.com/zetsux/gin-gorm-clean-starter/support/util"
	"gorm.io/gorm"
)

// --- Mock Repository ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) TxRepository() repository.TxRepository {
	args := m.Called()
	return args.Get(0).(repository.TxRepository)
}

func (m *MockUserRepository) CreateNewUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.User, error) {
	args := m.Called(ctx, tx, key, val)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context, tx *gorm.DB, req base.GetsRequest) ([]entity.User, int64, int64, error) {
	args := m.Called(ctx, tx, req)
	return args.Get(0).([]entity.User), args.Get(1).(int64), args.Get(2).(int64), args.Error(3)
}

func (m *MockUserRepository) UpdateNameUser(ctx context.Context, tx *gorm.DB, name string, user entity.User) (entity.User, error) {
	args := m.Called(ctx, tx, name, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user entity.User) (entity.User, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUserByID(ctx context.Context, tx *gorm.DB, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
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

func setupUserServiceMock() (service.UserService, *MockUserRepository, context.Context) {
	repo := new(MockUserRepository)
	us := service.NewUserService(repo)
	ctx := context.Background()

	return us, repo, ctx
}

// --- Tests ---

func TestUserService_CreateNewUser(t *testing.T) {
	us, repo, ctx := setupUserServiceMock()

	expected := entity.User{ID: uuid.New(), Name: "A", Email: "a@mail.test"}
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "email", "a@mail.test").Return(entity.User{}, errors.ErrUserNotFound)
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
	us, repo, ctx := setupUserServiceMock()

	hashed, _ := util.PasswordHash("secret")
	stored := entity.User{ID: uuid.New(), Email: "a@mail.test", Password: hashed}
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "email", "a@mail.test").Return(stored, nil)

	ok := us.VerifyLogin(ctx, "a@mail.test", "secret")
	require.True(t, ok)
	repo.AssertExpectations(t)
}

func TestUserService_GetAllUsers(t *testing.T) {
	us, repo, ctx := setupUserServiceMock()

	users := []entity.User{{ID: uuid.New(), Name: "A", Email: "a@mail.test"}}
	repo.On("GetAllUsers", ctx, (*gorm.DB)(nil), mock.AnythingOfType("base.GetsRequest")).Return(users, int64(1), int64(1), nil)

	got, page, err := us.GetAllUsers(ctx, base.GetsRequest{Page: 1, PerPage: 10})
	require.NoError(t, err)
	require.NotEmpty(t, got)
	require.Equal(t, int64(1), page.LastPage)
	repo.AssertExpectations(t)
}

func TestUserService_GetUserByPrimaryKey(t *testing.T) {
	us, repo, ctx := setupUserServiceMock()

	userID := uuid.New()
	expected := entity.User{ID: userID, Name: "A", Email: "a@mail.test"}
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "email", "a@mail.test").Return(expected, nil)

	fetched, err := us.GetUserByPrimaryKey(ctx, "email", "a@mail.test")
	require.NoError(t, err)
	require.Equal(t, userID.String(), fetched.ID)
	repo.AssertExpectations(t)
}

func TestUserService_UpdateSelfName(t *testing.T) {
	us, repo, ctx := setupUserServiceMock()

	userID := uuid.New()
	old := entity.User{ID: userID, Name: "A"}
	updated := entity.User{ID: userID, Name: "New"}

	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "id", userID.String()).Return(old, nil)
	repo.On("UpdateNameUser", ctx, (*gorm.DB)(nil), "New", old).Return(updated, nil)

	edited, err := us.UpdateSelfName(ctx, dto.UserNameUpdateRequest{Name: "New"}, userID.String())
	require.NoError(t, err)
	require.Equal(t, "New", edited.Name)
	repo.AssertExpectations(t)
}

func TestUserService_DeleteUserByID(t *testing.T) {
	us, repo, ctx := setupUserServiceMock()

	userID := uuid.New()
	repo.On("GetUserByPrimaryKey", ctx, (*gorm.DB)(nil), "id", userID.String()).Return(entity.User{ID: userID}, nil)
	repo.On("DeleteUserByID", ctx, (*gorm.DB)(nil), userID.String()).Return(nil)

	err := us.DeleteUserByID(ctx, userID.String())
	require.NoError(t, err)
	repo.AssertExpectations(t)
}
