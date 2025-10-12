package provider

import (
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/config"
	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"
	"gorm.io/gorm"
)

const DATABASE = "DATABASE"

func SetupDependencies(injector *do.Injector) {
	if _, err := do.InvokeNamed[*gorm.DB](injector, DATABASE); err != nil {
		do.ProvideNamed(injector, DATABASE, func(i *do.Injector) (*gorm.DB, error) {
			return config.DBSetup(), nil
		})
	}

	do.Provide(injector, func(i *do.Injector) (repository.TxRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, DATABASE)

		return repository.NewTxRepository(db), nil
	})

	do.Provide(injector, func(i *do.Injector) (service.JWTService, error) {
		return service.NewJWTService(), nil
	})

	SetupUserDependencies(injector)
	SetupFileDependencies(injector)
}
