package provider

import (
	"myapp/api/v1/controller"
	queryiface "myapp/core/interface/query"
	repositoryiface "myapp/core/interface/repository"
	"myapp/core/service"
	"myapp/infrastructure/query"
	"myapp/infrastructure/repository"
	"myapp/support/constant"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func SetupUserDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (repositoryiface.UserRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return repository.NewUserRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (queryiface.UserQuery, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return query.NewUserQuery(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.UserService, error) {
		userR := do.MustInvoke[repositoryiface.UserRepository](i)
		userQ := do.MustInvoke[queryiface.UserQuery](i)
		txR := do.MustInvoke[repositoryiface.TxRepository](i)
		return service.NewUserService(userR, userQ, txR), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.UserController, error) {
		userS := do.MustInvoke[service.UserService](i)
		jwtS := do.MustInvoke[service.JWTService](i)
		return controller.NewUserController(userS, jwtS), nil
	})
}
