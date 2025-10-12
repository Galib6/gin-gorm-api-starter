package provider

import (
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
)

func SetupUserDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (repository.UserRepository, error) {
		txR := do.MustInvoke[repository.TxRepository](i)
		return repository.NewUserRepository(txR), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.UserService, error) {
		userR := do.MustInvoke[repository.UserRepository](i)
		return service.NewUserService(userR), nil
	})

	do.Provide(injector, func(i *do.Injector) (controller.UserController, error) {
		userS := do.MustInvoke[service.UserService](i)
		jwtS := do.MustInvoke[service.JWTService](i)
		return controller.NewUserController(userS, jwtS), nil
	})
}
