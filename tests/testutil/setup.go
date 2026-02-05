package testutil

import (
	"testing"

	"myapp/api/v1/router"
	repositoryiface "myapp/core/interface/repository"
	"myapp/core/service"
	"myapp/provider"
	"myapp/support/constant"
	"myapp/support/middleware"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type TestApp struct {
	Server      *gin.Engine
	DB          *gorm.DB
	TxRepo      repositoryiface.TxRepository
	UserRepo    repositoryiface.UserRepository
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
	txRepo := do.MustInvoke[repositoryiface.TxRepository](injector)
	userRepo := do.MustInvoke[repositoryiface.UserRepository](injector)
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
