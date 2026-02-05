package service

import (
	"context"
	"fmt"
	"reflect"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	errs "myapp/core/helper/errors"
	queryiface "myapp/core/interface/query"
	repositoryiface "myapp/core/interface/repository"
	"myapp/support/base"
	"myapp/support/constant"
	"myapp/support/util"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type productService struct {
	productRepository  repositoryiface.ProductRepository
	categoryRepository repositoryiface.CategoryRepository
	productQuery       queryiface.ProductQuery
	categoryQuery      queryiface.CategoryQuery
	txRepository       repositoryiface.TxRepository
}

type ProductService interface {
	// Product CRUD
	CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (dto.ProductResponse, error)
	GetAllProducts(ctx context.Context, req dto.ProductGetsRequest) ([]dto.ProductResponse, base.PaginationResponse, error)
	GetProductByID(ctx context.Context, id string) (dto.ProductResponse, error)
	UpdateProduct(ctx context.Context, req dto.ProductUpdateRequest) (dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id string) error

	// Product Image
	ChangeProductImage(ctx context.Context, req dto.ProductChangeImageRequest) (dto.ProductResponse, error)
	DeleteProductImage(ctx context.Context, id string) error

	// Stock Management
	UpdateStock(ctx context.Context, req dto.ProductStockUpdateRequest) (dto.ProductResponse, error)

	// Complex Queries
	GetLowStockProducts(ctx context.Context, threshold int) ([]dto.ProductResponse, error)
	GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]dto.ProductResponse, error)
	GetProductStatsByCategory(ctx context.Context) ([]dto.CategoryProductStats, error)

	// Complex Maintenance Operation
	RunProductMaintenance(ctx context.Context, req dto.ProductMaintenanceRequest) (dto.ProductMaintenanceResponse, error)
}

func NewProductService(
	productR repositoryiface.ProductRepository,
	categoryR repositoryiface.CategoryRepository,
	productQ queryiface.ProductQuery,
	categoryQ queryiface.CategoryQuery,
	txR repositoryiface.TxRepository,
) ProductService {
	return &productService{
		productRepository:  productR,
		categoryRepository: categoryR,
		productQuery:       productQ,
		categoryQuery:      categoryQ,
		txRepository:       txR,
	}
}

// ============== Helper Functions ==============

func (sv *productService) toProductResponse(product entity.Product) dto.ProductResponse {
	resp := dto.ProductResponse{
		ID:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.SKU,
		Price:       product.Price,
		Stock:       product.Stock,
		IsActive:    product.IsActive,
	}

	if product.CategoryID != nil {
		resp.CategoryID = product.CategoryID.String()
	}

	if product.Image != nil {
		resp.Image = *product.Image
	}

	if product.Category != nil {
		resp.Category = &dto.CategoryResponse{
			ID:          product.Category.ID.String(),
			Name:        product.Category.Name,
			Description: product.Category.Description,
		}
	}

	return resp
}

// ============== Product CRUD ==============

func (sv *productService) CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (dto.ProductResponse, error) {
	// Check if SKU already exists
	existingProduct, err := sv.productRepository.GetProductByPrimaryKey(ctx, nil, constant.DBAttrSKU, req.SKU)
	if err != nil && err != errs.ErrProductNotFound {
		return dto.ProductResponse{}, err
	}
	if !reflect.DeepEqual(existingProduct, entity.Product{}) {
		return dto.ProductResponse{}, errs.ErrProductSKUExists
	}

	// Validate category if provided
	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		_, err := sv.categoryRepository.GetCategoryByID(ctx, nil, req.CategoryID)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		catUUID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		categoryID = &catUUID
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	product := entity.Product{
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       decimal.NewFromFloat(req.Price),
		Stock:       req.Stock,
		CategoryID:  categoryID,
		IsActive:    isActive,
	}

	newProduct, err := sv.productRepository.CreateProduct(ctx, nil, product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return sv.toProductResponse(newProduct), nil
}

func (sv *productService) GetAllProducts(ctx context.Context, req dto.ProductGetsRequest) (
	productsResp []dto.ProductResponse, pageResp base.PaginationResponse, err error) {

	products, pageResp, err := sv.productQuery.GetAllProducts(ctx, req)
	if err != nil {
		return []dto.ProductResponse{}, base.PaginationResponse{}, err
	}

	for _, product := range products {
		productsResp = append(productsResp, sv.toProductResponse(product))
	}
	return productsResp, pageResp, nil
}

func (sv *productService) GetProductByID(ctx context.Context, id string) (dto.ProductResponse, error) {
	product, err := sv.productRepository.GetProductByID(ctx, nil, id, "Category")
	if err != nil {
		return dto.ProductResponse{}, err
	}
	return sv.toProductResponse(product), nil
}

func (sv *productService) UpdateProduct(ctx context.Context, req dto.ProductUpdateRequest) (dto.ProductResponse, error) {
	product, err := sv.productRepository.GetProductByID(ctx, nil, req.ID)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Check if SKU is being changed and already exists
	if req.SKU != "" && req.SKU != product.SKU {
		existing, err := sv.productRepository.GetProductByPrimaryKey(ctx, nil, constant.DBAttrSKU, req.SKU)
		if err != nil && err != errs.ErrProductNotFound {
			return dto.ProductResponse{}, err
		}
		if !reflect.DeepEqual(existing, entity.Product{}) {
			return dto.ProductResponse{}, errs.ErrProductSKUExists
		}
	}

	// Validate new category if provided
	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		_, err := sv.categoryRepository.GetCategoryByID(ctx, nil, req.CategoryID)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		catUUID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			return dto.ProductResponse{}, err
		}
		categoryID = &catUUID
	}

	productEdit := entity.Product{
		ID:          product.ID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		CategoryID:  categoryID,
	}

	if req.Price != nil {
		productEdit.Price = decimal.NewFromFloat(*req.Price)
	}

	if req.Stock != nil {
		productEdit.Stock = *req.Stock
	}

	if req.IsActive != nil {
		productEdit.IsActive = *req.IsActive
	}

	err = sv.productRepository.UpdateProduct(ctx, nil, productEdit)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return sv.toProductResponse(productEdit), nil
}

func (sv *productService) DeleteProduct(ctx context.Context, id string) error {
	_, err := sv.productRepository.GetProductByID(ctx, nil, id)
	if err != nil {
		return err
	}

	return sv.productRepository.DeleteProductByID(ctx, nil, id)
}

// ============== Product Image ==============

func (sv *productService) ChangeProductImage(ctx context.Context, req dto.ProductChangeImageRequest) (dto.ProductResponse, error) {
	product, err := sv.productRepository.GetProductByID(ctx, nil, req.ID)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Delete old image if exists
	if product.Image != nil && *product.Image != "" {
		if err := util.DeleteFile(*product.Image); err != nil {
			return dto.ProductResponse{}, err
		}
	}

	// Upload new image
	imgID := uuid.New()
	imgPath := fmt.Sprintf("product_image/%v", imgID)
	if err := util.UploadFile(req.Image, imgPath); err != nil {
		return dto.ProductResponse{}, err
	}

	productEdit := entity.Product{
		ID:    product.ID,
		Image: &imgPath,
	}
	err = sv.productRepository.UpdateProduct(ctx, nil, productEdit)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return dto.ProductResponse{
		ID:    productEdit.ID.String(),
		Image: imgPath,
	}, nil
}

func (sv *productService) DeleteProductImage(ctx context.Context, id string) error {
	product, err := sv.productRepository.GetProductByID(ctx, nil, id)
	if err != nil {
		return err
	}

	if product.Image == nil || *product.Image == "" {
		return errs.ErrProductNoImage
	}

	if err := util.DeleteFile(*product.Image); err != nil {
		return err
	}

	emptyString := ""
	productEdit := entity.Product{
		ID:    product.ID,
		Image: &emptyString,
	}

	return sv.productRepository.UpdateProduct(ctx, nil, productEdit)
}

// ============== Stock Management ==============

func (sv *productService) UpdateStock(ctx context.Context, req dto.ProductStockUpdateRequest) (dto.ProductResponse, error) {
	product, err := sv.productRepository.GetProductByID(ctx, nil, req.ID)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Check if stock reduction would result in negative stock
	if req.Quantity < 0 && product.Stock+req.Quantity < 0 {
		return dto.ProductResponse{}, errs.ErrInsufficientStock
	}

	err = sv.productRepository.UpdateProductStock(ctx, nil, req.ID, req.Quantity)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	product.Stock += req.Quantity
	return sv.toProductResponse(product), nil
}

// ============== Complex Queries ==============

func (sv *productService) GetLowStockProducts(ctx context.Context, threshold int) ([]dto.ProductResponse, error) {
	products, err := sv.productQuery.GetLowStockProducts(ctx, threshold)
	if err != nil {
		return nil, err
	}

	var result []dto.ProductResponse
	for _, product := range products {
		result = append(result, sv.toProductResponse(product))
	}
	return result, nil
}

func (sv *productService) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]dto.ProductResponse, error) {
	if minPrice > maxPrice {
		return nil, errs.ErrInvalidPriceRange
	}

	products, err := sv.productQuery.GetProductsByPriceRange(ctx, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}

	var result []dto.ProductResponse
	for _, product := range products {
		result = append(result, sv.toProductResponse(product))
	}
	return result, nil
}

func (sv *productService) GetProductStatsByCategory(ctx context.Context) ([]dto.CategoryProductStats, error) {
	return sv.productQuery.GetProductStatsByCategory(ctx)
}

// ============== Complex Maintenance Operation ==============

// RunProductMaintenance is a complex batch operation that demonstrates:
//   - Using the query layer for complex, filterable queries
//   - Applying multiple business rules
//   - Performing multiple updates inside a single transaction
func (sv *productService) RunProductMaintenance(
	ctx context.Context,
	req dto.ProductMaintenanceRequest,
) (resp dto.ProductMaintenanceResponse, err error) {
	// 1. Normalize input / defaults
	if req.PerPage <= 0 || req.PerPage > 1000 {
		req.PerPage = 100
	}

	// 2. Build the underlying GetAllProducts request
	getReq := dto.ProductGetsRequest{
		ID:         req.ID,
		CategoryID: req.CategoryID,
		IsActive:   req.IsActive,
		Search:     req.Search,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		MinStock:   req.MinStock,
		MaxStock:   req.MaxStock,
		PaginationRequest: base.PaginationRequest{
			Sort:     req.Sort,
			Includes: req.Includes,
			Page:     req.Page,
			PerPage:  req.PerPage,
		},
	}

	// 3. Execute complex, filterable query via query layer
	products, pageResp, err := sv.productQuery.GetAllProducts(ctx, getReq)
	if err != nil {
		return dto.ProductMaintenanceResponse{}, err
	}
	resp.PaginationResponse = pageResp
	resp.TotalSelected = len(products)
	if len(products) == 0 {
		return resp, nil
	}

	// 4. Validate new category if provided
	if req.NewCategoryID != "" {
		_, err := sv.categoryRepository.GetCategoryByID(ctx, nil, req.NewCategoryID)
		if err != nil {
			return dto.ProductMaintenanceResponse{}, err
		}
	}

	// 5. Start transaction for all subsequent updates
	tx, err := sv.txRepository.BeginTx(ctx)
	if err != nil {
		return dto.ProductMaintenanceResponse{}, err
	}
	defer func() {
		sv.txRepository.CommitOrRollbackTx(ctx, tx, err)
	}()

	for _, product := range products {
		result := dto.ProductMaintenanceResult{
			ProductID:   product.ID.String(),
			ProductName: product.Name,
		}

		productEdit := entity.Product{
			ID: product.ID,
		}

		needsUpdate := false

		// Business rule: Adjust price if requested
		if req.PriceAdjustment != nil && *req.PriceAdjustment != 0 {
			oldPrice, _ := product.Price.Float64()
			multiplier := 1 + (*req.PriceAdjustment / 100)
			newPriceDecimal := product.Price.Mul(decimal.NewFromFloat(multiplier))
			newPrice, _ := newPriceDecimal.Float64()

			productEdit.Price = newPriceDecimal
			result.OldPrice = &oldPrice
			result.NewPrice = &newPrice
			resp.PriceChangedCount++
			needsUpdate = true
		}

		// Business rule: Set active status if requested
		if req.SetActive != nil && product.IsActive != *req.SetActive {
			productEdit.IsActive = *req.SetActive
			result.ActiveChanged = true
			resp.ActiveChangedCount++
			needsUpdate = true
		}

		// Business rule: Mark low stock products as inactive
		if req.LowStockThreshold != nil && product.Stock < *req.LowStockThreshold && product.IsActive {
			productEdit.IsActive = false
			result.ActiveChanged = true
			resp.ActiveChangedCount++
			needsUpdate = true
		}

		// Business rule: Move to new category if requested
		if req.NewCategoryID != "" {
			catUUID, parseErr := uuid.Parse(req.NewCategoryID)
			if parseErr != nil {
				err = parseErr
				return resp, err
			}
			if product.CategoryID == nil || product.CategoryID.String() != req.NewCategoryID {
				productEdit.CategoryID = &catUUID
				result.CategoryChanged = true
				resp.CategoryChangedCount++
				needsUpdate = true
			}
		}

		// Skip DB call if nothing changed
		if !needsUpdate {
			continue
		}

		if err = sv.productRepository.UpdateProduct(ctx, tx, productEdit); err != nil {
			return resp, err
		}

		resp.Details = append(resp.Details, result)
		resp.TotalProcessed++
	}

	return resp, nil
}
