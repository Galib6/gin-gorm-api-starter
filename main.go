package main

import (
	"fmt"
	"os"

	"myapp/api/v1/router"
	"myapp/cmd"
	"myapp/config"
	"myapp/database/migrations"
	docs "myapp/docs" // Swagger docs (allow runtime modifications)
	"myapp/provider"
	"myapp/support/constant"
	"myapp/support/logger"
	"myapp/support/middleware"

	"github.com/samber/do"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gin GORM API Starter
// @version         1.0
// @description     A starter API project with Gin and GORM featuring user management and product catalog.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize logger level based on environment
	if os.Getenv("DEBUG") == "true" {
		logger.SetLevel(logger.DEBUG)
	}

	logger.Info("🔧 Starting application...")

	// Setting Up Dependencies
	var injector = do.New()
	provider.SetupDependencies(injector)
	logger.Info("✅ Dependencies initialized")

	db := do.MustInvokeNamed[*gorm.DB](injector, constant.DBInjectorKey)
	defer config.DBClose(db)
	logger.Info("✅ Database connected")

	// Handling CLI Commands
	cmd.Execute(db)

	// Auto-run migrations on startup (if not executing a CLI command)
	if os.Getenv("SKIP_MIGRATION") != "true" {
		logger.Info("🔄 Checking for pending migrations...")
		if err := migrations.RunMigrations("up"); err != nil {
			logger.Error("❌ Failed to run auto-migrations: %v", err)
			// Don't exit, just log error. Or should we exit? Usually safe to exit in production if migration fails.
			// For dev/local, maybe we want to continue? Let's exit to be safe.
			os.Exit(1)
		}
		logger.Info("✅ Database is up to date")
	}

	// Setting Up Server with custom recovery and logger
	gin.SetMode(gin.ReleaseMode) // Disable default Gin logger
	server := gin.New()          // Use gin.New() instead of gin.Default() for custom middlewares

	// Apply middlewares in order:
	// 1. Recovery - catches panics and prevents crashes
	// 2. Request Logger - logs all incoming requests
	// 3. CORS - handles cross-origin requests
	// 4. Error Handler - handles errors set in context
	server.Use(
		middleware.RecoveryMiddleware(),
		middleware.RequestLoggerMiddleware(),
		middleware.CORSMiddleware(),
		middleware.ErrorHandler(),
	)

	// Swagger documentation endpoint
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setting Up Routes
	router.InitRoutes(server, injector)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set swagger host dynamically so Swagger UI reflects actual server port
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", port)

	fmt.Println()
	logger.Info("🚀 Server running on http://localhost:%s", port)
	logger.Info("📚 Swagger docs at http://localhost:%s/docs/index.html", port)
	fmt.Println()

	err := server.Run(":" + port)
	if err != nil {
		fmt.Println("Server failed to start: ", err)
		return
	}
}
