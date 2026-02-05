package provider

import (
	"myapp/config"
	repositoryiface "myapp/core/interface/repository"
	"myapp/core/service"
	"myapp/infrastructure/repository"
	"myapp/support/constant"

	"github.com/samber/do"
	"gorm.io/gorm"
)

func SetupDependencies(injector *do.Injector) {
	if _, err := do.InvokeNamed[*gorm.DB](injector, constant.DBInjectorKey); err != nil {
		do.ProvideNamed(injector, constant.DBInjectorKey, func(i *do.Injector) (*gorm.DB, error) {
			return config.DBSetup(), nil
		})
	}

	do.Provide(injector, func(i *do.Injector) (repositoryiface.TxRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)

		return repository.NewTxRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	SetupUserDependencies(injector)
	SetupFileDependencies(injector)
	SetupProductDependencies(injector)
}
