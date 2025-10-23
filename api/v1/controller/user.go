package controller

import (
	"net/http"
	"reflect"

	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-api-starter/core/service"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetAllUsers(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	UpdateSelfName(ctx *gin.Context)
	UpdateUserByID(ctx *gin.Context)
	DeleteSelfUser(ctx *gin.Context)
	DeleteUserByID(ctx *gin.Context)
	ChangePicture(ctx *gin.Context)
	DeletePicture(ctx *gin.Context)
}

func NewUserController(userS service.UserService, jwtS service.JWTService) UserController {
	return &userController{
		userService: userS,
		jwtService:  jwtS,
	}
}

func (ct *userController) Register(ctx *gin.Context) {
	var req dto.UserRegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUserRegisterFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	newUser, err := ct.userService.CreateNewUser(ctx, req)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserRegisterFailed, err))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		messages.MsgUserRegisterSuccess,
		http.StatusCreated, newUser,
	))
}

func (ct *userController) Login(ctx *gin.Context) {
	var req dto.UserLoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUserLoginFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	res := ct.userService.VerifyLogin(ctx, req.Email, req.Password)
	if !res {
		_ = ctx.Error(base.NewAppError(http.StatusUnauthorized,
			messages.MsgUserWrongCredential, nil))
		return
	}

	user, err := ct.userService.GetUserByPrimaryKey(ctx, constant.DBAttrEmail, req.Email)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserLoginFailed, err))
		return
	}

	token := ct.jwtService.GenerateToken(user.ID, user.Role)
	authResp := base.CreateAuthResponse(token, user.Role)
	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserLoginSuccess,
		http.StatusOK, authResp,
	))
}

func (ct *userController) GetAllUsers(ctx *gin.Context) {
	var req dto.UserGetsRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUsersFetchFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	users, pageMeta, err := ct.userService.GetAllUsers(ctx, req)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUsersFetchFailed, err))
		return
	}

	if reflect.DeepEqual(pageMeta, base.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
			messages.MsgUsersFetchSuccess,
			http.StatusOK, users,
		))
	} else {
		ctx.JSON(http.StatusOK, base.CreatePaginatedResponse(
			messages.MsgUsersFetchSuccess,
			http.StatusOK, users, pageMeta,
		))
	}
}

func (ct *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := ct.userService.GetUserByPrimaryKey(ctx, constant.DBAttrID, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserFetchFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserFetchSuccess,
		http.StatusOK, user,
	))
}

func (ct *userController) UpdateSelfName(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var req dto.UserNameUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUserUpdateFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	req.ID = id
	user, err := ct.userService.UpdateSelfName(ctx, req)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserUpdateFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (ct *userController) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")

	var req dto.UserUpdateRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUserUpdateFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	req.ID = id
	user, err := ct.userService.UpdateUserByID(ctx, req)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserUpdateFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserUpdateSuccess,
		http.StatusOK, user,
	))
}

func (ct *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	err := ct.userService.DeleteUserByID(ctx, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserDeleteFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (ct *userController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := ct.userService.DeleteUserByID(ctx, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserDeleteFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserDeleteSuccess,
		http.StatusOK, nil,
	))
}

func (ct *userController) ChangePicture(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)

	var req dto.UserChangePictureRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUserPictureUpdateFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	res, err := ct.userService.ChangePicture(ctx, req, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserPictureUpdateFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserPictureUpdateSuccess,
		http.StatusOK, res,
	))
}

func (ct *userController) DeletePicture(ctx *gin.Context) {
	id := ctx.Param("user_id")
	err := ct.userService.DeletePicture(ctx, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserPictureDeleteFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserPictureDeleteSuccess,
		http.StatusOK, nil,
	))
}
