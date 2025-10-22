package router

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-api-starter/api/v1/controller"
)

func FileRouter(route *gin.Engine, injector *do.Injector) {
	var (
		fileController = do.MustInvoke[controller.FileController](injector)
	)

	routes := route.Group("/api/v1/files")
	{
		routes.GET("/:dir/:file_id", fileController.GetFile)
	}
}
