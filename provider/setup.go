package provider

import (
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/config"
	repository_interface "github.com/zetsux/gin-gorm-clean-starter/core/interface/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"github.com/zetsux/gin-gorm-clean-starter/infrastructure/repository"
	"github.com/zetsux/gin-gorm-clean-starter/support/constant"
	"gorm.io/gorm"
)

func SetupDependencies(injector *do.Injector) {
	if _, err := do.InvokeNamed[*gorm.DB](injector, constant.DBInjectorKey); err != nil {
		do.ProvideNamed(injector, constant.DBInjectorKey, func(i *do.Injector) (*gorm.DB, error) {
			return config.DBSetup(), nil
		})
	}

	do.Provide(injector, func(i *do.Injector) (repository_interface.TxRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)

		return repository.NewTxRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	SetupUserDependencies(injector)
	SetupFileDependencies(injector)
}
