package controller

import (
	"net/http"
	"os"
	"strings"

	"github.com/zetsux/gin-gorm-api-starter/core/helper/messages"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"

	"github.com/gin-gonic/gin"
)

type fileController struct{}

type FileController interface {
	GetFile(ctx *gin.Context)
}

func NewFileController() FileController {
	return &fileController{}
}

func (ct *fileController) GetFile(ctx *gin.Context) {
	dir := ctx.Param("dir")
	fileID := ctx.Param("file_id")

	filePath := strings.Join([]string{constant.FileBasePath, dir, fileID}, "/")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgFileFetchFailed, err))
		return
	}

	ctx.File(filePath)
}
