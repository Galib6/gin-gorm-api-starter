package provider

import (
	"github.com/samber/do"
	"github.com/zetsux/gin-gorm-clean-starter/api/v1/controller"
)

func SetupFileDependencies(injector *do.Injector) {
	do.Provide(injector, func(i *do.Injector) (controller.FileController, error) {
		return controller.NewFileController(), nil
	})
}
