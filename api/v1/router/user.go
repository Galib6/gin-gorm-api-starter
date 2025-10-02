package router

import (
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/common/middleware"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, userC controller.UserController, jwtS service.JWTService) {
	userRoutes := router.Group("/api/v1/users")
	{
		// admin routes
		userRoutes.GET("", middleware.Authenticate(jwtS), middleware.Authorize(), userC.GetAllUsers)
		userRoutes.PATCH("/:user_id", middleware.Authenticate(jwtS), middleware.Authorize(), userC.UpdateUserByID)
		userRoutes.DELETE("/:user_id", middleware.Authenticate(jwtS), middleware.Authorize(), userC.DeleteUserByID)

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
