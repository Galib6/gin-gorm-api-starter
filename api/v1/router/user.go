package router

import (
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-api-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-api-starter/core/service"
	"github.com/zetsux/gin-gorm-api-starter/support/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, injector *do.Injector) {
	var (
		userC = do.MustInvoke[controller.UserController](injector)
		jwtS  = do.MustInvoke[service.JWTService](injector)
	)

	userRoutes := router.Group("/api/v1/users")
	{
		// admin routes
		userRoutes.GET("", middleware.Authenticate(jwtS), middleware.Authorize(), userC.GetAllUsers)
		userRoutes.PATCH("/:user_id", middleware.Authenticate(jwtS), middleware.Authorize(), userC.UpdateUserByID)
		userRoutes.DELETE("/:user_id", middleware.Authenticate(jwtS), middleware.Authorize(), userC.DeleteUserByID)
		userRoutes.POST("/maintenance", middleware.Authenticate(jwtS), middleware.Authorize(), userC.RunUserMaintenance)

		// user routes
		userRoutes.GET("/me", middleware.Authenticate(jwtS), userC.GetMe)
		userRoutes.PATCH("/me/name", middleware.Authenticate(jwtS), userC.UpdateSelfName)
		userRoutes.DELETE("/me", middleware.Authenticate(jwtS), userC.DeleteSelfUser)
		userRoutes.POST("", userC.Register)
		userRoutes.POST("/login", userC.Login)

		// user file routes
		userRoutes.PATCH("/picture", middleware.Authenticate(jwtS), userC.ChangePicture)
		userRoutes.DELETE("/picture/:user_id", middleware.Authenticate(jwtS), userC.DeletePicture)
	}
}
