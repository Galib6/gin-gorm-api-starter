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

func SetupProductDependencies(injector *do.Injector) {
	// Product Repository
	do.Provide(injector, func(i *do.Injector) (repositoryiface.ProductRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return repository.NewProductRepository(db), nil
	})

	// Category Repository
	do.Provide(injector, func(i *do.Injector) (repositoryiface.CategoryRepository, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return repository.NewCategoryRepository(db), nil
	})

	// Product Query
	do.Provide(injector, func(i *do.Injector) (queryiface.ProductQuery, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return query.NewProductQuery(db), nil
	})

	// Category Query
	do.Provide(injector, func(i *do.Injector) (queryiface.CategoryQuery, error) {
		db := do.MustInvokeNamed[*gorm.DB](i, constant.DBInjectorKey)
		return query.NewCategoryQuery(db), nil
	})

	// Product Service
	do.Provide(injector, func(i *do.Injector) (service.ProductService, error) {
		productR := do.MustInvoke[repositoryiface.ProductRepository](i)
		categoryR := do.MustInvoke[repositoryiface.CategoryRepository](i)
		productQ := do.MustInvoke[queryiface.ProductQuery](i)
		categoryQ := do.MustInvoke[queryiface.CategoryQuery](i)
		txR := do.MustInvoke[repositoryiface.TxRepository](i)
		return service.NewProductService(productR, categoryR, productQ, categoryQ, txR), nil
	})

	// Category Service
	do.Provide(injector, func(i *do.Injector) (service.CategoryService, error) {
		categoryR := do.MustInvoke[repositoryiface.CategoryRepository](i)
		productR := do.MustInvoke[repositoryiface.ProductRepository](i)
		categoryQ := do.MustInvoke[queryiface.CategoryQuery](i)
		productQ := do.MustInvoke[queryiface.ProductQuery](i)
		return service.NewCategoryService(categoryR, productR, categoryQ, productQ), nil
	})

	// Product Controller
	do.Provide(injector, func(i *do.Injector) (controller.ProductController, error) {
		productS := do.MustInvoke[service.ProductService](i)
		categoryS := do.MustInvoke[service.CategoryService](i)
		return controller.NewProductController(productS, categoryS), nil
	})
}
