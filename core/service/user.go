package service

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	errs "myapp/core/helper/errors"
	queryiface "myapp/core/interface/query"
	repositoryiface "myapp/core/interface/repository"
	"myapp/support/base"
	"myapp/support/constant"
	"myapp/support/util"

	"github.com/google/uuid"
)

type userService struct {
	userRepository repositoryiface.UserRepository
	userQuery      queryiface.UserQuery
	txRepository   repositoryiface.TxRepository
}

type UserService interface {
	VerifyLogin(ctx context.Context, email string, password string) bool
	CreateNewUser(ctx context.Context, req dto.UserRegisterRequest) (dto.UserResponse, error)
	GetAllUsers(ctx context.Context, req dto.UserGetsRequest) ([]dto.UserResponse, base.PaginationResponse, error)
	GetUserByPrimaryKey(ctx context.Context, key string, value string) (dto.UserResponse, error)
	UpdateSelfName(ctx context.Context, req dto.UserNameUpdateRequest) (dto.UserResponse, error)
	UpdateUserByID(ctx context.Context, req dto.UserUpdateRequest) (dto.UserResponse, error)
	DeleteUserByID(ctx context.Context, id string) error
	ChangePicture(ctx context.Context, req dto.UserChangePictureRequest) (dto.UserResponse, error)
	DeletePicture(ctx context.Context, userID string) error
	// RunUserMaintenance demonstrates a complex, highly-customizable operation
	// that can apply multiple business rules and database changes in one
	// transactional flow.
	RunUserMaintenance(ctx context.Context, req dto.UserMaintenanceRequest) (dto.UserMaintenanceResponse, error)
}

func NewUserService(userR repositoryiface.UserRepository, userQ queryiface.UserQuery,
	txR repositoryiface.TxRepository,
) UserService {
	return &userService{
		userRepository: userR,
		userQuery:      userQ,
		txRepository:   txR,
	}
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
	_, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, id)
	if err != nil {
		return err
	}

	err = sv.userRepository.DeleteUserByID(ctx, nil, id)
	if err != nil {
		return err
	}
	return nil
}

func (sv *userService) ChangePicture(ctx context.Context,
	req dto.UserChangePictureRequest) (dto.UserResponse, error) {
	user, err := sv.userRepository.GetUserByPrimaryKey(ctx, nil, constant.DBAttrID, req.ID)
	if err != nil {
		return dto.UserResponse{}, err
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

// RunUserMaintenance is a "large" example function that shows how to:
//   - run complex, customizable queries via the query layer
//   - apply multiple business rules
//   - perform multiple updates inside a single transaction
func (sv *userService) RunUserMaintenance(
	ctx context.Context,
	req dto.UserMaintenanceRequest,
) (resp dto.UserMaintenanceResponse, err error) {
	// 1. Normalize input / defaults
	if req.PerPage <= 0 || req.PerPage > 1000 {
		req.PerPage = 100
	}
	if req.InactiveDays < 0 {
		req.InactiveDays = 0
	}

	// 2. Build the underlying GetAllUsers request from the embedded struct
	getReq := dto.UserGetsRequest{
		ID:     req.ID,
		Role:   req.Role,
		Search: req.Search,
		PaginationRequest: base.PaginationRequest{
			Sort:     req.Sort,
			Includes: req.Includes,
			Page:     req.Page,
			PerPage:  req.PerPage,
		},
	}

	// 3. Execute complex, filterable query via query layer
	users, pageResp, err := sv.userQuery.GetAllUsers(ctx, getReq)
	if err != nil {
		return dto.UserMaintenanceResponse{}, err
	}
	resp.PaginationResponse = pageResp
	resp.TotalSelected = len(users)
	if len(users) == 0 {
		return resp, nil
	}

	// 4. Start transaction for all subsequent updates
	tx, err := sv.txRepository.BeginTx(ctx)
	if err != nil {
		return dto.UserMaintenanceResponse{}, err
	}
	defer func() {
		sv.txRepository.CommitOrRollbackTx(ctx, tx, err)
	}()

	now := time.Now()

	for _, user := range users {
		// Optional inactivity filter based on UpdatedAt
		if req.InactiveDays > 0 {
			inactiveForDays := int(now.Sub(user.UpdatedAt).Hours() / 24)
			if inactiveForDays < req.InactiveDays {
				// Skip users that are not inactive long enough
				continue
			}
		}

		result := dto.UserMaintenanceResult{
			UserID:  user.ID.String(),
			OldRole: user.Role,
		}

		// Construct a minimal edit model
		userEdit := entity.User{
			ID: user.ID,
		}

		// Business rule: change role if requested
		if req.NewRole != "" && user.Role != req.NewRole {
			userEdit.Role = req.NewRole
			result.NewRole = req.NewRole
			resp.RoleChangedCount++
		} else {
			userEdit.Role = user.Role
		}

		// Business rule: clear picture if requested
		if req.ClearPicture && user.Picture != nil && *user.Picture != "" {
			// Remove file from storage
			if errDel := util.DeleteFile(*user.Picture); errDel != nil {
				err = errDel
				return resp, err
			}
			empty := ""
			userEdit.Picture = &empty
			result.PictureCleared = true
			resp.PictureClearedCount++
		} else {
			userEdit.Picture = user.Picture
		}

		// If nothing changed for this user, skip DB call
		if result.NewRole == "" && !result.PictureCleared {
			continue
		}

		if err = sv.userRepository.UpdateUser(ctx, tx, userEdit); err != nil {
			return resp, err
		}

		resp.Details = append(resp.Details, result)
		resp.TotalProcessed++
	}

	return resp, nil
}
