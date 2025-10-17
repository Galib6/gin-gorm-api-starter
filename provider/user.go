package provider

import (
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	query_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/query"
	repository_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/query"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-clean-starter/support/constant"
	"gorm.io/gorm"
)

func SetupUserDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (repository_interface.UserRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return repository.NewUserRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (query_interface.UserQuery, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return query.NewUserQuery(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.UserService, error) {
		userR := do.MustInvoke[repository_interface.UserRepository](i)
		userQ := do.MustInvoke[query_interface.UserQuery](i)
		return service.NewUserService(userR, userQ), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.UserController, error) {
		userS := do.MustInvoke[service.UserService](i)
		jwtS := do.MustInvoke[service.JWTService](i)
		return controller.NewUserController(userS, jwtS), nil
	})
}
