package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/router"
	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
)

type userServiceMock struct{ mock.Mock }

func (m *userServiceMock) VerifyLogin(ctx context.Context, email string, password string) bool {
	args := m.Called(ctx, email, password)
	return args.Bool(0)
}
func (m *userServiceMock) CreateNewUser(ctx context.Context, ud dto.UserRegisterRequest) (dto.UserResponse, error) {
	args := m.Called(ctx, ud)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) GetAllUsers(ctx context.Context, req base.GetsRequest) ([]dto.UserResponse, base.PaginationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]dto.UserResponse), args.Get(1).(base.PaginationResponse), args.Error(2)
}
func (m *userServiceMock) GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) UpdateSelfName(ctx context.Context, ud dto.UserNameUpdateRequest, id string) (dto.UserResponse, error) {
	args := m.Called(ctx, ud, id)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) UpdateUserByID(ctx context.Context, ud dto.UserUpdateRequest, id string) (dto.UserResponse, error) {
	args := m.Called(ctx, ud, id)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) DeleteUserByID(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *userServiceMock) ChangePicture(ctx context.Context, req dto.UserChangePictureRequest, userID string) (dto.UserResponse, error) {
	args := m.Called(ctx, req, userID)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) DeletePicture(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type jwtServiceMock struct{ mock.Mock }

func (j *jwtServiceMock) GenerateToken(id string, role string) string { return "token" }
func (j *jwtServiceMock) ValidateToken(token string) (*jwt.Token, error) {
	return &jwt.Token{Valid: true}, nil
}
func (j *jwtServiceMock) GetAttrByToken(token string) (string, string, error) {
	return "id", "user", nil
}

func TestUserController_RegisterAndLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	usm := new(userServiceMock)
	jwtm := new(jwtServiceMock)

	userC := controller.NewUserController(usm, jwtm)
	router.UserRouter(r, userC, jwtm)

	// Register
	regReq := dto.UserRegisterRequest{Name: "A", Email: "a@mail.test", Password: "secret"}
	usm.On("CreateNewUser", mock.Anything, regReq).Return(dto.UserResponse{ID: "1", Email: regReq.Email, Name: regReq.Name}, nil)
	b, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Login
	loginReq := dto.UserLoginRequest{Email: regReq.Email, Password: regReq.Password}
	usm.On("VerifyLogin", mock.Anything, loginReq.Email, loginReq.Password).Return(true)
	usm.On("GetUserByPrimaryKey", mock.Anything, "email", loginReq.Email).Return(dto.UserResponse{ID: "1", Email: loginReq.Email, Name: "A", Role: "user"}, nil)
	b, _ = json.Marshal(loginReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
