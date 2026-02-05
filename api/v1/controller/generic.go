package controller

import (
	"context"
	"net/http"
	"reflect"

	"myapp/support/base"

	"github.com/gin-gonic/gin"
)

func HandleCreate[T any, R any](
	ctx *gin.Context,
	dto T,
	createFunc func(context.Context, T) (R, error),
	successMsg, failMsg string,
) {
	if err := ctx.ShouldBind(&dto); err != nil {
		msg := base.GetValidationErrorMessage(err, dto, failMsg)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	result, err := createFunc(ctx, dto)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, failMsg, err))
		return
	}

	ctx.JSON(http.StatusCreated, base.CreateSuccessResponse(
		successMsg, http.StatusCreated, result))
}

func HandleGetAll[T any, R any](
	ctx *gin.Context,
	dto T,
	getAllFunc func(context.Context, T) (R, base.PaginationResponse, error),
	successMsg, failMsg string,
) {
	if err := ctx.ShouldBind(&dto); err != nil {
		msg := base.GetValidationErrorMessage(err, dto, failMsg)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	results, pageMeta, err := getAllFunc(ctx, dto)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, failMsg, err))
		return
	}

	if reflect.DeepEqual(pageMeta, base.PaginationResponse{}) {
		ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
			successMsg, http.StatusOK, results,
		))
	} else {
		ctx.JSON(http.StatusOK, base.CreatePaginatedResponse(
			successMsg, http.StatusOK, results, pageMeta,
		))
	}
}

func HandleGetByID[R any](
	ctx *gin.Context,
	id string,
	getByIDFunc func(context.Context, string) (R, error),
	successMsg, failMsg string,
) {
	result, err := getByIDFunc(ctx, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, failMsg, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		successMsg, http.StatusOK, result,
	))
}

func HandleUpdate[T any, R any](
	ctx *gin.Context,
	id string,
	dto T,
	updateFunc func(context.Context, T) (R, error),
	successMsg, failMsg string,
) {
	if err := ctx.ShouldBind(&dto); err != nil {
		msg := base.GetValidationErrorMessage(err, dto, failMsg)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	v := reflect.ValueOf(&dto).Elem()
	if field := v.FieldByName("ID"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(id)
	}

	result, err := updateFunc(ctx, dto)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, failMsg, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		successMsg, http.StatusOK, result,
	))
}

func HandleDelete(
	ctx *gin.Context,
	id string,
	deleteFunc func(context.Context, string) error,
	successMsg, failMsg string,
) {
	err := deleteFunc(ctx, id)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, failMsg, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		successMsg, http.StatusOK, nil,
	))
}
