package testutil

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/router"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/provider"
	"github.com/zetsux/gin-gorm-clean-starter/support/middleware"
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

	// Setup Dependencies
	injector := do.New()
	do.ProvideNamed(injector, provider.DATABASE, func(i *do.Injector) (*gorm.DB, error) {
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
	testDB := do.MustInvokeNamed[*gorm.DB](injector, provider.DATABASE)
	txRepo := do.MustInvoke[repository.TxRepository](injector)
	userRepo := do.MustInvoke[repository.UserRepository](injector)
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
