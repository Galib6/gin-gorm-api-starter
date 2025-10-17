package testutil

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/router"
	repository_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/provider"
	"github.com/zetsux/gin-gorm-clean-starter/support/constant"
	"github.com/zetsux/gin-gorm-clean-starter/support/middleware"
	"gorm.io/gorm"
)

type TestApp struct {
	Server      *gin.Engine
	DB          *gorm.DB
	TxRepo      repository_interface.TxRepository
	UserRepo    repository_interface.UserRepository
	UserService service.UserService
	JWTService  service.JWTService
}

func SetupTestApp(t *testing.T) *TestApp {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// Setup Dependencies
	injector := do.New()
	do.ProvideNamed(injector, constant.DBInjectorKey, func(i *do.Injector) (*gorm.DB, error) {
		return NewTestDB(t), nil
	})
	provider.SetupDependencies(injector)

	// Router
	r := gin.New()
	r.Use(
		middleware.CORSMiddleware(),
		middleware.ErrorHandler(),
	)
	router.UserRouter(r, injector)

	// Invoke
	testDB := do.MustInvokeNamed[*gorm.DB](injector, constant.DBInjectorKey)
	txRepo := do.MustInvoke[repository_interface.TxRepository](injector)
	userRepo := do.MustInvoke[repository_interface.UserRepository](injector)
	userService := do.MustInvoke[service.UserService](injector)
	jwtService := do.MustInvoke[service.JWTService](injector)

	return &TestApp{
		Server:      r,
		DB:          testDB,
		TxRepo:      txRepo,
		UserRepo:    userRepo,
		UserService: userService,
		JWTService:  jwtService,
	}
}
