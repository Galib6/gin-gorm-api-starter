package router

import (
	"myapp/api/v1/controller"
	"myapp/core/service"
	"myapp/support/middleware"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"
)

func ProductRouter(router *gin.Engine, injector *do.Injector) {
	var (
		productC = do.MustInvoke[controller.ProductController](injector)
		jwtS     = do.MustInvoke[service.JWTService](injector)
	)

	// ============== Product Routes ==============
	productRoutes := router.Group("/api/v1/products")
	{
		// Public routes
		productRoutes.GET("", productC.GetAllProducts)
		productRoutes.GET("/:product_id", productC.GetProductByID)
		productRoutes.GET("/low-stock", productC.GetLowStockProducts)
		productRoutes.GET("/price-range", productC.GetProductsByPriceRange)
		productRoutes.GET("/stats/by-category", productC.GetProductStatsByCategory)

		// Admin routes (require authentication and authorization)
		productRoutes.POST("", middleware.Authenticate(jwtS), middleware.Authorize(), productC.CreateProduct)
		productRoutes.PATCH("/:product_id", middleware.Authenticate(jwtS), middleware.Authorize(), productC.UpdateProduct)
		productRoutes.DELETE("/:product_id", middleware.Authenticate(jwtS), middleware.Authorize(), productC.DeleteProduct)

		// Product image routes
		productRoutes.PATCH("/:product_id/image", middleware.Authenticate(jwtS), middleware.Authorize(), productC.ChangeProductImage)
		productRoutes.DELETE("/:product_id/image", middleware.Authenticate(jwtS), middleware.Authorize(), productC.DeleteProductImage)

		// Stock management routes
		productRoutes.PATCH("/:product_id/stock", middleware.Authenticate(jwtS), middleware.Authorize(), productC.UpdateStock)

		// Complex maintenance operation
		productRoutes.POST("/maintenance", middleware.Authenticate(jwtS), middleware.Authorize(), productC.RunProductMaintenance)
	}

	// ============== Category Routes ==============
	categoryRoutes := router.Group("/api/v1/categories")
	{
		// Public routes
		categoryRoutes.GET("", productC.GetAllCategories)
		categoryRoutes.GET("/:category_id", productC.GetCategoryByID)

		// Admin routes
		categoryRoutes.POST("", middleware.Authenticate(jwtS), middleware.Authorize(), productC.CreateCategory)
		categoryRoutes.PATCH("/:category_id", middleware.Authenticate(jwtS), middleware.Authorize(), productC.UpdateCategory)
		categoryRoutes.DELETE("/:category_id", middleware.Authenticate(jwtS), middleware.Authorize(), productC.DeleteCategory)
	}
}
