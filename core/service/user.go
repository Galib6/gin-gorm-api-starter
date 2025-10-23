package service

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/zetsux/gin-gorm-api-starter/core/entity"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	errs "github.com/zetsux/gin-gorm-api-starter/core/helper/errors"
	queryiface "github.com/zetsux/gin-gorm-api-starter/core/interface/query"
	repositoryiface "github.com/zetsux/gin-gorm-api-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"
	"github.com/zetsux/gin-gorm-api-starter/support/util"
)

type userService struct {
	userRepository repositoryiface.UserRepository
	userQuery      queryiface.UserQuery
}

type UserService interface {
	VerifyLogin(ctx context.Context, email string, password string) bool
	CreateNewUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context, req dto.UserGetsRequest) ([]dto.UserResponse, base.PaginationResponse, error)
	GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error)
	UpdateSelfName(ctx context.Context, req dto.UserNameUpdateRequest) (dto.UserResponse, error)
	UpdateUserByID(ctx context.Context, req dto.UserUpdateRequest) (dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, id string) error
	ChangePicture(ctx context.Context, req dto.UserChangePictureRequest, userID string) (dto.UserResponse, error)
	DeletePicture(ctx context.Context, userID string) error
}

func NewUserService(userR repositoryiface.UserRepository, userQ queryiface.UserQuery) UserService {
	return &userService{userRepository: userR, userQuery: userQ}
}

func (sv *userService) VerifyLogin(ctx context.Context, email string, password string) bool {
	userCheck, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrEmail, email)
	if err != nil {
		return false
	}
	passwordCheck, err := util.PasswordCompare(userCheck.Password, []byte(password))
	if err != nil {
		return false
	}

	if userCheck.Email == email && passwordCheck {
		return true
	}
	return false
}

func (sv *userService) CreateNewUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserResponse, error) {
	userCheck, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrEmail, req.Email)
	if err != nil && err != errs.ErrUserNotFound {
		return dto.UserResponse{}, err
	}

	if !(reflect.DeepEqual(userCheck, entity.User{})) {
		return dto.UserResponse{}, errs.ErrEmailAlreadyExists
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     constant.EnumRoleUser,
	}

	// create new user
	newUser, err := sv.userRepository.CreateNewUser(ctx, nil, user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    newUser.ID.String(),
		Name:  newUser.Name,
		Email: newUser.Email,
		Role:  newUser.Role,
	}, nil
}

func (sv *userService) GetAllUsers(ctx context.Context, req dto.UserGetsRequest) (
	usersResp []dto.UserResponse, pageResp base.PaginationResponse, err error) {
	users, pageResp, err := sv.userQuery.GetAllUsers(ctx, req)
	if err != nil {
		return []dto.UserResponse{}, base.PaginationResponse{}, err
	}

	for _, user := range users {
		userResp := dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}
		if user.Picture != nil {
			userResp.Picture = *user.Picture
		}

		usersResp = append(usersResp, userResp)
	}
	return usersResp, pageResp, nil
}

func (sv *userService) GetUserByPrimaryKey(ctx context.Context, key string, val string) (dto.UserResponse, error) {
	user, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, key, val)
	if err != nil {
		return dto.UserResponse{}, err
	}

	userResp := dto.UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	if user.Picture != nil {
		userResp.Picture = *user.Picture
	}

	return userResp, nil
}

func (sv *userService) UpdateSelfName(ctx context.Context,
	req dto.UserNameUpdateRequest) (dto.UserResponse, error) {
	user, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, req.ID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	userEdit := entity.User{
		ID:   user.ID,
		Name: req.Name,
	}
	err = sv.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:   userEdit.ID.String(),
		Name: userEdit.Name,
	}, nil
}

func (sv *userService) UpdateUserByID(ctx context.Context,
	req dto.UserUpdateRequest) (dto.UserResponse, error) {
	user, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, req.ID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, errs.ErrUserNotFound
	}

	if req.Email != "" && req.Email != user.Email {
		us, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrEmail, req.Email)
		if err != nil && err != errs.ErrUserNotFound {
			return dto.UserResponse{}, err
		}

		if !(reflect.DeepEqual(us, entity.User{})) {
			return dto.UserResponse{}, errs.ErrEmailAlreadyExists
		}
	}

	userEdit := entity.User{
		ID:       user.ID,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	}
	err = sv.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		ID:    userEdit.ID.String(),
		Name:  userEdit.Name,
		Email: userEdit.Email,
		Role:  userEdit.Role,
	}, nil
}

func (sv *userService) DeleteUserByID(ctx context.Context, id string) error {
	userCheck, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, id)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(userCheck, entity.User{}) {
		return errs.ErrUserNotFound
	}

	err = sv.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}

func (sv *userService) ChangePicture(ctx context.Context,
	req dto.UserChangePictureRequest, userID string) (dto.UserResponse, error) {
	user, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return dto.UserResponse{}, errs.ErrUserNotFound
	}

	if user.Picture != nil && *user.Picture != "" {
		if err := util.DeleteFile(*user.Picture); err != nil {
			return dto.UserResponse{}, err
		}
	}

	picID := uuid.New()
	picPath := fmt.Sprintf("user_picture/%v", picID)
	if err := util.UploadFile(req.Picture, picPath); err != nil {
		return dto.UserResponse{}, err
	}

	userEdit := entity.User{
		ID:      user.ID,
		Picture: &picPath,
	}
	err = sv.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return dto.UserResponse{}, err
	}

	userResp := dto.UserResponse{
		ID:      userEdit.ID.String(),
		Picture: picPath,
	}
	return userResp, nil
}

func (sv *userService) DeletePicture(ctx context.Context, userID string) error {
	user, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, userID)
	if err != nil {
		return err
	}

	if reflect.DeepEqual(user, entity.User{}) {
		return errs.ErrUserNotFound
	}

	if user.Picture == nil || *user.Picture == "" {
		return errs.ErrUserNoPicture
	}

	if err := util.DeleteFile(*user.Picture); err != nil {
		return err
	}

	emptyString := ""
	userEdit := entity.User{
		ID:      user.ID,
		Picture: &emptyString,
	}

	err = sv.userRepository.UpdateUser(ctx, nil, userEdit)
	if err != nil {
		return err
	}

	return nil
}
