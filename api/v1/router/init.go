package router

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func InitRoutes(server *gin.Engine, injector *do.Injector) {
	UserRouter(server, injector)
	FileRouter(server, injector)
	ProductRouter(server, injector)
}
