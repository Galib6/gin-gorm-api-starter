package provider

import (
	"myapp/api/v1/controller"

	"github.com/samber/do"
)

func SetupFileDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (controller.FileController, error) {
		return controller.NewFileController(), nil
	})
}
