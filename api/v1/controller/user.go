package controller

import (
	"net/http"

	"myapp/core/helper/dto"
	"myapp/core/helper/messages"
	"myapp/core/service"
	"myapp/support/base"
	"myapp/support/constant"

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
	// RunUserMaintenance is an example of a large, complex operation that
	// applies many business rules and DB changes in one call.
	RunUserMaintenance(ctx *gin.Context)
}

func NewUserController(userS service.UserService, jwtS service.JWTService) UserController {
	return &userController{
		userService: userS,
		jwtService:  jwtS,
	}
}

func (uc *userController) Register(ctx *gin.Context) {
	HandleCreate(ctx, dto.UserRegisterRequest{}, uc.userService.CreateNewUser,
		messages.MsgUserRegisterSuccess, messages.MsgUserRegisterFailed)
}

func (uc *userController) Login(ctx *gin.Context) {
	var userDTO dto.UserLoginRequest
	if err := ctx.ShouldBind(&userDTO); err != nil {
		msg := base.GetValidationErrorMessage(err, userDTO, messages.MsgUserLoginFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	res := uc.userService.VerifyLogin(ctx, userDTO.Email, userDTO.Password)
	if !res {
		_ = ctx.Error(base.NewAppError(http.StatusUnauthorized,
			messages.MsgUserWrongCredential, nil))
		return
	}

	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrEmail, userDTO.Email)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgUserLoginFailed, err))
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	authResp := base.CreateAuthResponse(token, user.Role)
	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserLoginSuccess,
		http.StatusOK, authResp,
	))
}

func (uc *userController) GetAllUsers(ctx *gin.Context) {
	HandleGetAll(ctx, dto.UserGetsRequest{}, uc.userService.GetAllUsers,
		messages.MsgUsersFetchSuccess, messages.MsgUsersFetchFailed)
}

func (uc *userController) GetMe(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	user, err := uc.userService.GetUserByPrimaryKey(ctx, constant.DBAttrID, id)
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

func (uc *userController) UpdateSelfName(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	HandleUpdate(ctx, id, dto.UserNameUpdateRequest{}, uc.userService.UpdateSelfName,
		messages.MsgUserUpdateSuccess, messages.MsgUserUpdateFailed)
}

func (uc *userController) UpdateUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	HandleUpdate(ctx, id, dto.UserUpdateRequest{}, uc.userService.UpdateUserByID,
		messages.MsgUserUpdateSuccess, messages.MsgUserUpdateFailed)
}

func (uc *userController) DeleteSelfUser(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	HandleDelete(ctx, id, uc.userService.DeleteUserByID,
		messages.MsgUserDeleteSuccess, messages.MsgUserDeleteFailed)
}

func (uc *userController) DeleteUserByID(ctx *gin.Context) {
	id := ctx.Param("user_id")
	HandleDelete(ctx, id, uc.userService.DeleteUserByID,
		messages.MsgUserDeleteSuccess, messages.MsgUserDeleteFailed)
}

func (uc *userController) ChangePicture(ctx *gin.Context) {
	id := ctx.MustGet("ID").(string)
	HandleUpdate(ctx, id, dto.UserChangePictureRequest{}, uc.userService.ChangePicture,
		messages.MsgUserPictureUpdateSuccess, messages.MsgUserPictureUpdateFailed)
}

func (uc *userController) DeletePicture(ctx *gin.Context) {
	id := ctx.Param("user_id")
	HandleDelete(ctx, id, uc.userService.DeletePicture,
		messages.MsgUserPictureDeleteSuccess, messages.MsgUserPictureDeleteFailed)
}

// RunUserMaintenance demonstrates how to call a "large" service function from
// the controller. It binds a complex request DTO and returns a summary
// response.
func (uc *userController) RunUserMaintenance(ctx *gin.Context) {
	var req dto.UserMaintenanceRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgUserUpdateFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	resp, err := uc.userService.RunUserMaintenance(ctx, req)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusInternalServerError,
			messages.MsgUserUpdateFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgUserUpdateSuccess,
		http.StatusOK,
		resp,
	))
}
