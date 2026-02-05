package controller

import (
	"net/http"
	"strconv"

	"myapp/core/helper/dto"
	"myapp/core/helper/messages"
	"myapp/core/service"
	"myapp/support/base"

	"github.com/gin-gonic/gin"
)

type productController struct {
	productService  service.ProductService
	categoryService service.CategoryService
}

type ProductController interface {
	// Product CRUD
	CreateProduct(ctx *gin.Context)
	GetAllProducts(ctx *gin.Context)
	GetProductByID(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)

	// Product Image
	ChangeProductImage(ctx *gin.Context)
	DeleteProductImage(ctx *gin.Context)

	// Stock Management
	UpdateStock(ctx *gin.Context)

	// Complex Queries
	GetLowStockProducts(ctx *gin.Context)
	GetProductsByPriceRange(ctx *gin.Context)
	GetProductStatsByCategory(ctx *gin.Context)

	// Complex Maintenance
	RunProductMaintenance(ctx *gin.Context)

	// Category CRUD
	CreateCategory(ctx *gin.Context)
	GetAllCategories(ctx *gin.Context)
	GetCategoryByID(ctx *gin.Context)
	UpdateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
}

func NewProductController(productS service.ProductService, categoryS service.CategoryService) ProductController {
	return &productController{
		productService:  productS,
		categoryService: categoryS,
	}
}

// ============== Product CRUD ==============

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with the provided details
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product  body      dto.ProductCreateRequest  true  "Product details"
// @Success      201      {object}  base.Response{data=dto.ProductResponse}
// @Failure      400      {object}  base.Response
// @Security     BearerAuth
// @Router       /products [post]
func (pc *productController) CreateProduct(ctx *gin.Context) {
	HandleCreate(ctx, dto.ProductCreateRequest{}, pc.productService.CreateProduct,
		messages.MsgProductCreateSuccess, messages.MsgProductCreateFailed)
}

// GetAllProducts godoc
// @Summary      Get all products
// @Description  Get all products with filtering, sorting, and pagination
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        filter[id]          query     string  false  "Filter by product ID"
// @Param        filter[category_id] query     string  false  "Filter by category ID"
// @Param        filter[is_active]   query     bool    false  "Filter by active status"
// @Param        filter[min_price]   query     number  false  "Filter by minimum price"
// @Param        filter[max_price]   query     number  false  "Filter by maximum price"
// @Param        filter[min_stock]   query     int     false  "Filter by minimum stock"
// @Param        filter[max_stock]   query     int     false  "Filter by maximum stock"
// @Param        search              query     string  false  "Search in name, description, SKU"
// @Param        sort                query     string  false  "Sort field (prefix with - for desc)"
// @Param        page                query     int     false  "Page number"
// @Param        per_page            query     int     false  "Items per page"
// @Param        includes            query     string  false  "Include relations (e.g., Category)"
// @Success      200                 {object}  base.Response{data=[]dto.ProductResponse}
// @Failure      400                 {object}  base.Response
// @Router       /products [get]
func (pc *productController) GetAllProducts(ctx *gin.Context) {
	HandleGetAll(ctx, dto.ProductGetsRequest{}, pc.productService.GetAllProducts,
		messages.MsgProductsFetchSuccess, messages.MsgProductsFetchFailed)
}

// GetProductByID godoc
// @Summary      Get product by ID
// @Description  Get a single product by its ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id  path      string  true  "Product ID"
// @Success      200         {object}  base.Response{data=dto.ProductResponse}
// @Failure      400         {object}  base.Response
// @Router       /products/{product_id} [get]
func (pc *productController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("product_id")
	HandleGetByID(ctx, id, pc.productService.GetProductByID,
		messages.MsgProductFetchSuccess, messages.MsgProductFetchFailed)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update product details by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id  path      string                    true  "Product ID"
// @Param        product     body      dto.ProductUpdateRequest  true  "Product update details"
// @Success      200         {object}  base.Response{data=dto.ProductResponse}
// @Failure      400         {object}  base.Response
// @Security     BearerAuth
// @Router       /products/{product_id} [patch]
func (pc *productController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("product_id")
	HandleUpdate(ctx, id, dto.ProductUpdateRequest{}, pc.productService.UpdateProduct,
		messages.MsgProductUpdateSuccess, messages.MsgProductUpdateFailed)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id  path      string  true  "Product ID"
// @Success      200         {object}  base.Response
// @Failure      400         {object}  base.Response
// @Security     BearerAuth
// @Router       /products/{product_id} [delete]
func (pc *productController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("product_id")
	HandleDelete(ctx, id, pc.productService.DeleteProduct,
		messages.MsgProductDeleteSuccess, messages.MsgProductDeleteFailed)
}

// ============== Product Image ==============

// ChangeProductImage godoc
// @Summary      Upload product image
// @Description  Upload or change product image
// @Tags         Products
// @Accept       multipart/form-data
// @Produce      json
// @Param        product_id  path      string  true  "Product ID"
// @Param        image       formData  file    true  "Product image"
// @Success      200         {object}  base.Response{data=dto.ProductResponse}
// @Failure      400         {object}  base.Response
// @Security     BearerAuth
// @Router       /products/{product_id}/image [patch]
func (pc *productController) ChangeProductImage(ctx *gin.Context) {
	id := ctx.Param("product_id")
	HandleUpdate(ctx, id, dto.ProductChangeImageRequest{}, pc.productService.ChangeProductImage,
		messages.MsgProductImageUpdateSuccess, messages.MsgProductImageUpdateFailed)
}

// DeleteProductImage godoc
// @Summary      Delete product image
// @Description  Delete the image of a product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id  path      string  true  "Product ID"
// @Success      200         {object}  base.Response
// @Failure      400         {object}  base.Response
// @Security     BearerAuth
// @Router       /products/{product_id}/image [delete]
func (pc *productController) DeleteProductImage(ctx *gin.Context) {
	id := ctx.Param("product_id")
	HandleDelete(ctx, id, pc.productService.DeleteProductImage,
		messages.MsgProductImageDeleteSuccess, messages.MsgProductImageDeleteFailed)
}

// ============== Stock Management ==============

// UpdateStock godoc
// @Summary      Update product stock
// @Description  Add or subtract stock quantity
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_id  path      string                        true  "Product ID"
// @Param        stock       body      dto.ProductStockUpdateRequest true  "Stock update (positive to add, negative to subtract)"
// @Success      200         {object}  base.Response{data=dto.ProductResponse}
// @Failure      400         {object}  base.Response
// @Security     BearerAuth
// @Router       /products/{product_id}/stock [patch]
func (pc *productController) UpdateStock(ctx *gin.Context) {
	id := ctx.Param("product_id")
	HandleUpdate(ctx, id, dto.ProductStockUpdateRequest{}, pc.productService.UpdateStock,
		messages.MsgProductStockUpdateSuccess, messages.MsgProductStockUpdateFailed)
}

// ============== Complex Queries ==============

// GetLowStockProducts godoc
// @Summary      Get low stock products
// @Description  Get products with stock below the specified threshold
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        threshold  query     int  false  "Stock threshold (default: 10)"
// @Success      200        {object}  base.Response{data=[]dto.ProductResponse}
// @Failure      400        {object}  base.Response
// @Router       /products/low-stock [get]
func (pc *productController) GetLowStockProducts(ctx *gin.Context) {
	thresholdStr := ctx.DefaultQuery("threshold", "10")
	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgProductsFetchFailed, err))
		return
	}

	products, err := pc.productService.GetLowStockProducts(ctx, threshold)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgProductsFetchFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgProductsFetchSuccess,
		http.StatusOK, products,
	))
}

// GetProductsByPriceRange godoc
// @Summary      Get products by price range
// @Description  Get all products within a specific price range
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        min_price  query     number  true  "Minimum price"
// @Param        max_price  query     number  true  "Maximum price"
// @Success      200        {object}  base.Response{data=[]dto.ProductResponse}
// @Failure      400        {object}  base.Response
// @Router       /products/price-range [get]
func (pc *productController) GetProductsByPriceRange(ctx *gin.Context) {
	minPriceStr := ctx.Query("min_price")
	maxPriceStr := ctx.Query("max_price")

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgProductsFetchFailed, err))
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgProductsFetchFailed, err))
		return
	}

	products, err := pc.productService.GetProductsByPriceRange(ctx, minPrice, maxPrice)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgProductsFetchFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgProductsFetchSuccess,
		http.StatusOK, products,
	))
}

// GetProductStatsByCategory godoc
// @Summary      Get product statistics by category
// @Description  Get aggregated product statistics grouped by category
// @Tags         Products
// @Accept       json
// @Produce      json
// @Success      200  {object}  base.Response{data=[]dto.CategoryProductStats}
// @Failure      400  {object}  base.Response
// @Router       /products/stats/by-category [get]
func (pc *productController) GetProductStatsByCategory(ctx *gin.Context) {
	stats, err := pc.productService.GetProductStatsByCategory(ctx)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest,
			messages.MsgProductStatsFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgProductStatsSuccess,
		http.StatusOK, stats,
	))
}

// ============== Complex Maintenance ==============

// RunProductMaintenance godoc
// @Summary      Run product maintenance
// @Description  Batch operation to update multiple products based on filters
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        maintenance  body      dto.ProductMaintenanceRequest  true  "Maintenance parameters"
// @Success      200          {object}  base.Response{data=dto.ProductMaintenanceResponse}
// @Failure      400          {object}  base.Response
// @Failure      500          {object}  base.Response
// @Security     BearerAuth
// @Router       /products/maintenance [post]
func (pc *productController) RunProductMaintenance(ctx *gin.Context) {
	var req dto.ProductMaintenanceRequest
	if err := ctx.ShouldBind(&req); err != nil {
		msg := base.GetValidationErrorMessage(err, req, messages.MsgProductMaintenanceFailed)
		_ = ctx.Error(base.NewAppError(http.StatusBadRequest, msg, err))
		return
	}

	resp, err := pc.productService.RunProductMaintenance(ctx, req)
	if err != nil {
		_ = ctx.Error(base.NewAppError(http.StatusInternalServerError,
			messages.MsgProductMaintenanceFailed, err))
		return
	}

	ctx.JSON(http.StatusOK, base.CreateSuccessResponse(
		messages.MsgProductMaintenanceSuccess,
		http.StatusOK,
		resp,
	))
}

// ============== Category CRUD ==============

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Create a new product category
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category  body      dto.CategoryCreateRequest  true  "Category details"
// @Success      201       {object}  base.Response{data=dto.CategoryResponse}
// @Failure      400       {object}  base.Response
// @Security     BearerAuth
// @Router       /categories [post]
func (pc *productController) CreateCategory(ctx *gin.Context) {
	HandleCreate(ctx, dto.CategoryCreateRequest{}, pc.categoryService.CreateCategory,
		messages.MsgCategoryCreateSuccess, messages.MsgCategoryCreateFailed)
}

// GetAllCategories godoc
// @Summary      Get all categories
// @Description  Get all categories with optional filtering and pagination
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        filter[id]  query     string  false  "Filter by category ID"
// @Param        search      query     string  false  "Search in name and description"
// @Param        sort        query     string  false  "Sort field (prefix with - for desc)"
// @Param        page        query     int     false  "Page number"
// @Param        per_page    query     int     false  "Items per page"
// @Success      200         {object}  base.Response{data=[]dto.CategoryResponse}
// @Failure      400         {object}  base.Response
// @Router       /categories [get]
func (pc *productController) GetAllCategories(ctx *gin.Context) {
	HandleGetAll(ctx, dto.CategoryGetsRequest{}, pc.categoryService.GetAllCategories,
		messages.MsgCategoriesFetchSuccess, messages.MsgCategoriesFetchFailed)
}

// GetCategoryByID godoc
// @Summary      Get category by ID
// @Description  Get a single category by its ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category_id  path      string  true  "Category ID"
// @Success      200          {object}  base.Response{data=dto.CategoryResponse}
// @Failure      400          {object}  base.Response
// @Router       /categories/{category_id} [get]
func (pc *productController) GetCategoryByID(ctx *gin.Context) {
	id := ctx.Param("category_id")
	HandleGetByID(ctx, id, pc.categoryService.GetCategoryByID,
		messages.MsgCategoryFetchSuccess, messages.MsgCategoryFetchFailed)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Update category details by ID
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category_id  path      string                     true  "Category ID"
// @Param        category     body      dto.CategoryUpdateRequest  true  "Category update details"
// @Success      200          {object}  base.Response{data=dto.CategoryResponse}
// @Failure      400          {object}  base.Response
// @Security     BearerAuth
// @Router       /categories/{category_id} [patch]
func (pc *productController) UpdateCategory(ctx *gin.Context) {
	id := ctx.Param("category_id")
	HandleUpdate(ctx, id, dto.CategoryUpdateRequest{}, pc.categoryService.UpdateCategory,
		messages.MsgCategoryUpdateSuccess, messages.MsgCategoryUpdateFailed)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Delete a category by ID (fails if category has products)
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        category_id  path      string  true  "Category ID"
// @Success      200          {object}  base.Response
// @Failure      400          {object}  base.Response
// @Security     BearerAuth
// @Router       /categories/{category_id} [delete]
func (pc *productController) DeleteCategory(ctx *gin.Context) {
	id := ctx.Param("category_id")
	HandleDelete(ctx, id, pc.categoryService.DeleteCategory,
		messages.MsgCategoryDeleteSuccess, messages.MsgCategoryDeleteFailed)
}
