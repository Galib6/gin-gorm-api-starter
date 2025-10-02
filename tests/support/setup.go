package support

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/router"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"gorm.io/gorm"
)

type TestApp struct {
	Server      *gin.Engine
	DB          *gorm.DB
	TxRepo      repository.TxRepository
	UserRepo    repository.UserRepository
	UserService service.UserService
	JWTService  service.JWTService
}

func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// DB
	db := NewTestDB(t)

	// Repos
	txRepo := repository.NewTxRepository(db)
	userRepo := repository.NewUserRepository(txRepo)

	// Services
	userService := service.NewUserService(userRepo)
	jwtService := service.NewJWTService()

	// Controllers
	userController := controller.NewUserController(userService, jwtService)

	// Router
	r := gin.New()
	router.UserRouter(r, userController, jwtService)

	return &TestApp{
		Server:      r,
		DB:          db,
		TxRepo:      txRepo,
		UserRepo:    userRepo,
		UserService: userService,
		JWTService:  jwtService,
	}
}
