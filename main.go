package main

import (
	"fmt"
	"os"

	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-api-starter/api/v1/router"
	"github.com/zetsux/gin-gorm-api-starter/config"
	"github.com/zetsux/gin-gorm-api-starter/provider"
	"github.com/zetsux/gin-gorm-api-starter/support/constant"
	"github.com/zetsux/gin-gorm-api-starter/support/middleware"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setting Up Dependencies
	var injector = do.New()
	provider.SetupDependencies(injector)

	db := do.MustInvokeNamed[*gorm.DB](injector, constant.DBInjectorKey)
	defer config.DBClose(db)

	// Setting Up Server
	server := gin.Default()
	server.Use(
		middleware.CORSMiddleware(),
		middleware.ErrorHandler(),
	)

	// Setting Up Routes
	router.InitRoutes(server, injector)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := server.Run(":" + port)
	if err != nil {
		fmt.Println("Server failed to start: ", err)
		return
	}
}
