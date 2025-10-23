package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zetsux/gin-gorm-api-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-api-starter/api/v1/router"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-api-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-api-starter/core/service"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"
	"github.com/zetsux/gin-gorm-api-starter/support/middleware"
)

// --- Mock Services ---

type userServiceMock struct{ mock.Mock }

func (m *userServiceMock) VerifyLogin(ctx context.Context, email string, password string) bool {
	args := m.Called(ctx, email, password)
	return args.Bool(0)
}
func (m *userServiceMock) CreateNewUser(ctx context.Context, ud dto.UserRegisterRequest) (dto.UserResponse, error) {
	args := m.Called(ctx, ud)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) GetAllUsers(ctx context.Context, req dto.UserGetsRequest) ([]dto.UserResponse, base.PaginationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]dto.UserResponse), args.Get(1).(base.PaginationResponse), args.Error(2)
}
func (m *userServiceMock) GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) UpdateSelfName(ctx context.Context, ud dto.UserNameUpdateRequest) (dto.UserResponse, error) {
	args := m.Called(ctx, ud)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
func (m *userServiceMock) UpdateUserByID(ctx context.Context, ud dto.UserUpdateRequest) (dto.UserResponse, error) {
	args := m.Called(ctx, ud)
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
	args := j.Called(token)
	return args.String(0), args.String(1), args.Error(2)
}

// --- Test Helpers ---

func setupUserControllerTest() (*gin.Engine, *userServiceMock, *jwtServiceMock) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	// Setup dependencies
	injector := do.New()
	usm := new(userServiceMock)
	jwtm := new(jwtServiceMock)
	userC := controller.NewUserController(usm, jwtm)
	do.Provide(injector, func(i *do.Injector) (service.JWTService, error) {
		return jwtm, nil
	})
	do.Provide(injector, func(i *do.Injector) (controller.UserController, error) {
		return userC, nil
	})

	router.UserRouter(r, injector)
	return r, usm, jwtm
}

// --- Tests ---

func TestUserController_Register(t *testing.T) {
	r, usm, _ := setupUserControllerTest()

	regReq := dto.UserRegisterRequest{Name: "A", Email: "a@mail.test", Password: "secret"}
	usm.On("CreateNewUser", mock.Anything, regReq).Return(dto.UserResponse{ID: uuid.NewString(), Email: regReq.Email, Name: regReq.Name}, nil)

	b, _ := json.Marshal(regReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}

func TestUserController_Login(t *testing.T) {
	r, usm, _ := setupUserControllerTest()

	loginReq := dto.UserLoginRequest{Email: "a@mail.test", Password: "secret"}
	usm.On("VerifyLogin", mock.Anything, loginReq.Email, loginReq.Password).Return(true)
	usm.On("GetUserByPrimaryKey", mock.Anything, "email", loginReq.Email).Return(
		dto.UserResponse{ID: uuid.NewString(), Email: loginReq.Email, Name: "A", Role: constant.EnumRoleUser}, nil,
	)

	b, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_Login_Invalid(t *testing.T) {
	r, usm, _ := setupUserControllerTest()

	loginReq := dto.UserLoginRequest{Email: "a@mail.test", Password: "wrong"}
	usm.On("VerifyLogin", mock.Anything, loginReq.Email, loginReq.Password).Return(false)
	usm.On("GetUserByPrimaryKey", mock.Anything, "email", loginReq.Email).Return(
		dto.UserResponse{}, errs.ErrUserNotFound,
	)

	b, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserController_GetMe(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	uuidStr := uuid.NewString()
	jwtm.On("GetAttrByToken", "token").Return(uuidStr, constant.EnumRoleUser, nil)
	usm.On("GetUserByPrimaryKey", mock.Anything, mock.Anything, mock.Anything).Return(
		dto.UserResponse{ID: uuidStr, Email: "a@mail.test", Name: "A", Role: constant.EnumRoleUser}, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_GetMe_Unauthenticated(t *testing.T) {
	r, _, _ := setupUserControllerTest()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserController_GetAllUsers(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	getsReq := dto.UserGetsRequest{Search: "a", PaginationRequest: base.PaginationRequest{Page: 1, PerPage: 10}}
	jwtm.On("GetAttrByToken", "token").Return(uuid.NewString(), constant.EnumRoleAdmin, nil)
	usm.On("GetAllUsers", mock.Anything, getsReq).Return(
		[]dto.UserResponse{
			{ID: uuid.NewString(), Email: "a@mail.test", Name: "A", Role: constant.EnumRoleUser},
		}, base.PaginationResponse{Page: 1, PerPage: 10, Total: 1}, nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?search=a&page=1&per_page=10", nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_UpdateSelfName(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	uuidStr := uuid.NewString()
	jwtm.On("GetAttrByToken", "token").Return(uuidStr, constant.EnumRoleUser, nil)

	updateReq := dto.UserNameUpdateRequest{ID: uuidStr, Name: "Updated Name"}
	usm.On("UpdateSelfName", mock.Anything, updateReq).Return(
		dto.UserResponse{ID: uuidStr, Email: "a@mail.test", Name: "Updated Name", Role: constant.EnumRoleUser}, nil,
	)

	b, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/me/name", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_UpdateUserByID(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	targetUserID := uuid.NewString()
	jwtm.On("GetAttrByToken", "token").Return(uuid.NewString(), constant.EnumRoleAdmin, nil)

	updateReq := dto.UserUpdateRequest{ID: targetUserID, Name: "Updated Name", Role: constant.EnumRoleUser}
	usm.On("UpdateUserByID", mock.Anything, updateReq).Return(
		dto.UserResponse{ID: targetUserID, Email: "a@mail.test", Name: "Updated Name", Role: constant.EnumRoleUser}, nil,
	)

	b, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/"+targetUserID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_Delete(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	targetUserID := uuid.NewString()
	jwtm.On("GetAttrByToken", "token").Return(uuid.NewString(), constant.EnumRoleAdmin, nil)
	usm.On("DeleteUserByID", mock.Anything, targetUserID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+targetUserID, nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_ChangePicture(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	uuidStr := uuid.NewString()
	jwtm.On("GetAttrByToken", "token").Return(uuidStr, constant.EnumRoleUser, nil)

	// Create a multipart form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("picture", "profile.png")
	require.NoError(t, err)

	_, err = io.Copy(fileWriter, strings.NewReader("fake image bytes"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/users/picture", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	err = req.ParseMultipartForm(10 << 20)
	require.NoError(t, err)
	fileHeader := req.MultipartForm.File["picture"][0]

	changePicReq := dto.UserChangePictureRequest{Picture: fileHeader}
	usm.On("ChangePicture", mock.Anything, changePicReq, uuidStr).Return(
		dto.UserResponse{ID: uuidStr}, nil,
	)

	usm.On("ChangePicture", mock.Anything, changePicReq, uuidStr).
		Return(dto.UserResponse{ID: uuidStr, Name: "Updated"}, nil)

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUserController_DeletePicture(t *testing.T) {
	r, usm, jwtm := setupUserControllerTest()

	uuidStr := uuid.NewString()
	jwtm.On("GetAttrByToken", "token").Return(uuidStr, constant.EnumRoleUser, nil)

	usm.On("DeletePicture", mock.Anything, uuidStr).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/picture/"+uuidStr, nil)
	req.Header.Set("Authorization", "Bearer token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
