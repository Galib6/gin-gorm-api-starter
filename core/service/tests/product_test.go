package service

import (
	"context"
	"testing"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	errs "myapp/core/helper/errors"
	"myapp/core/service"
	"myapp/support/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// ============== Mock Repositories ==============

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) DB() *gorm.DB {
	return nil
}

func (m *mockProductRepository) CreateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) (entity.Product, error) {
	args := m.Called(ctx, tx, product)
	return args.Get(0).(entity.Product), args.Error(1)
}

func (m *mockProductRepository) GetProductByID(ctx context.Context, tx *gorm.DB, id string, includes ...string) (entity.Product, error) {
	args := m.Called(ctx, tx, id, includes)
	return args.Get(0).(entity.Product), args.Error(1)
}

func (m *mockProductRepository) GetProductByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.Product, error) {
	args := m.Called(ctx, tx, key, val)
	return args.Get(0).(entity.Product), args.Error(1)
}

func (m *mockProductRepository) UpdateProduct(ctx context.Context, tx *gorm.DB, product entity.Product) error {
	args := m.Called(ctx, tx, product)
	return args.Error(0)
}

func (m *mockProductRepository) DeleteProductByID(ctx context.Context, tx *gorm.DB, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

func (m *mockProductRepository) UpdateProductStock(ctx context.Context, tx *gorm.DB, id string, quantity int) error {
	args := m.Called(ctx, tx, id, quantity)
	return args.Error(0)
}

func (m *mockProductRepository) BulkUpdatePrices(ctx context.Context, tx *gorm.DB, ids []string, priceMultiplier float64) error {
	args := m.Called(ctx, tx, ids, priceMultiplier)
	return args.Error(0)
}

type mockCategoryRepository struct {
	mock.Mock
}

func (m *mockCategoryRepository) DB() *gorm.DB {
	return nil
}

func (m *mockCategoryRepository) CreateCategory(ctx context.Context, tx *gorm.DB, category entity.Category) (entity.Category, error) {
	args := m.Called(ctx, tx, category)
	return args.Get(0).(entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) GetCategoryByID(ctx context.Context, tx *gorm.DB, id string, includes ...string) (entity.Category, error) {
	args := m.Called(ctx, tx, id, includes)
	return args.Get(0).(entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) GetCategoryByPrimaryKey(ctx context.Context, tx *gorm.DB, key string, val string) (entity.Category, error) {
	args := m.Called(ctx, tx, key, val)
	return args.Get(0).(entity.Category), args.Error(1)
}

func (m *mockCategoryRepository) UpdateCategory(ctx context.Context, tx *gorm.DB, category entity.Category) error {
	args := m.Called(ctx, tx, category)
	return args.Error(0)
}

func (m *mockCategoryRepository) DeleteCategoryByID(ctx context.Context, tx *gorm.DB, id string) error {
	args := m.Called(ctx, tx, id)
	return args.Error(0)
}

// ============== Mock Queries ==============

type mockProductQuery struct {
	mock.Mock
}

func (m *mockProductQuery) GetAllProducts(ctx context.Context, req dto.ProductGetsRequest) ([]entity.Product, base.PaginationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]entity.Product), args.Get(1).(base.PaginationResponse), args.Error(2)
}

func (m *mockProductQuery) GetProductsByPriceRange(ctx context.Context, minPrice, maxPrice float64) ([]entity.Product, error) {
	args := m.Called(ctx, minPrice, maxPrice)
	return args.Get(0).([]entity.Product), args.Error(1)
}

func (m *mockProductQuery) GetLowStockProducts(ctx context.Context, threshold int) ([]entity.Product, error) {
	args := m.Called(ctx, threshold)
	return args.Get(0).([]entity.Product), args.Error(1)
}

func (m *mockProductQuery) GetProductStatsByCategory(ctx context.Context) ([]dto.CategoryProductStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.CategoryProductStats), args.Error(1)
}

type mockCategoryQuery struct {
	mock.Mock
}

func (m *mockCategoryQuery) GetAllCategories(ctx context.Context, req dto.CategoryGetsRequest) ([]entity.Category, base.PaginationResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]entity.Category), args.Get(1).(base.PaginationResponse), args.Error(2)
}

type mockTxRepository struct {
	mock.Mock
}

func (m *mockTxRepository) DB() *gorm.DB {
	return nil
}

func (m *mockTxRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *mockTxRepository) CommitOrRollbackTx(ctx context.Context, tx *gorm.DB, err error) {
	m.Called(ctx, tx, err)
}

// ============== Tests ==============

func TestCreateProduct_Success(t *testing.T) {
	// Setup
	mockProductRepo := new(mockProductRepository)
	mockCategoryRepo := new(mockCategoryRepository)
	mockProductQ := new(mockProductQuery)
	mockCategoryQ := new(mockCategoryQuery)
	mockTxRepo := new(mockTxRepository)

	productService := service.NewProductService(
		mockProductRepo, mockCategoryRepo, mockProductQ, mockCategoryQ, mockTxRepo,
	)

	ctx := context.Background()
	productID := uuid.New()

	req := dto.ProductCreateRequest{
		Name:        "Test Product",
		Description: "Test Description",
		SKU:         "TEST-SKU-001",
		Price:       99.99,
		Stock:       100,
	}

	expectedProduct := entity.Product{
		ID:          productID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Price:       decimal.NewFromFloat(req.Price),
		Stock:       req.Stock,
		IsActive:    true,
	}

	// Expectations
	mockProductRepo.On("GetProductByPrimaryKey", ctx, (*gorm.DB)(nil), "sku", req.SKU).
		Return(entity.Product{}, errs.ErrProductNotFound)
	mockProductRepo.On("CreateProduct", ctx, (*gorm.DB)(nil), mock.AnythingOfType("entity.Product")).
		Return(expectedProduct, nil)

	// Execute
	result, err := productService.CreateProduct(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.Equal(t, expectedProduct.SKU, result.SKU)
	mockProductRepo.AssertExpectations(t)
}

func TestCreateProduct_DuplicateSKU(t *testing.T) {
	// Setup
	mockProductRepo := new(mockProductRepository)
	mockCategoryRepo := new(mockCategoryRepository)
	mockProductQ := new(mockProductQuery)
	mockCategoryQ := new(mockCategoryQuery)
	mockTxRepo := new(mockTxRepository)

	productService := service.NewProductService(
		mockProductRepo, mockCategoryRepo, mockProductQ, mockCategoryQ, mockTxRepo,
	)

	ctx := context.Background()

	req := dto.ProductCreateRequest{
		Name:  "Test Product",
		SKU:   "EXISTING-SKU",
		Price: 99.99,
	}

	existingProduct := entity.Product{
		ID:   uuid.New(),
		Name: "Existing Product",
		SKU:  "EXISTING-SKU",
	}

	// Expectations
	mockProductRepo.On("GetProductByPrimaryKey", ctx, (*gorm.DB)(nil), "sku", req.SKU).
		Return(existingProduct, nil)

	// Execute
	_, err := productService.CreateProduct(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errs.ErrProductSKUExists, err)
	mockProductRepo.AssertExpectations(t)
}

func TestGetProductByID_Success(t *testing.T) {
	// Setup
	mockProductRepo := new(mockProductRepository)
	mockCategoryRepo := new(mockCategoryRepository)
	mockProductQ := new(mockProductQuery)
	mockCategoryQ := new(mockCategoryQuery)
	mockTxRepo := new(mockTxRepository)

	productService := service.NewProductService(
		mockProductRepo, mockCategoryRepo, mockProductQ, mockCategoryQ, mockTxRepo,
	)

	ctx := context.Background()
	productID := uuid.New()
	categoryID := uuid.New()

	expectedProduct := entity.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		SKU:         "TEST-SKU-001",
		Price:       decimal.NewFromFloat(99.99),
		Stock:       100,
		CategoryID:  &categoryID,
		IsActive:    true,
		Category: &entity.Category{
			ID:   categoryID,
			Name: "Test Category",
		},
	}

	// Expectations
	mockProductRepo.On("GetProductByID", ctx, (*gorm.DB)(nil), productID.String(), []string{"Category"}).
		Return(expectedProduct, nil)

	// Execute
	result, err := productService.GetProductByID(ctx, productID.String())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID.String(), result.ID)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.NotNil(t, result.Category)
	mockProductRepo.AssertExpectations(t)
}

func TestUpdateStock_Success(t *testing.T) {
	// Setup
	mockProductRepo := new(mockProductRepository)
	mockCategoryRepo := new(mockCategoryRepository)
	mockProductQ := new(mockProductQuery)
	mockCategoryQ := new(mockCategoryQuery)
	mockTxRepo := new(mockTxRepository)

	productService := service.NewProductService(
		mockProductRepo, mockCategoryRepo, mockProductQ, mockCategoryQ, mockTxRepo,
	)

	ctx := context.Background()
	productID := uuid.New()

	existingProduct := entity.Product{
		ID:    productID,
		Name:  "Test Product",
		Stock: 100,
	}

	req := dto.ProductStockUpdateRequest{
		ID:       productID.String(),
		Quantity: 50,
	}

	// Expectations
	mockProductRepo.On("GetProductByID", ctx, (*gorm.DB)(nil), productID.String(), []string(nil)).
		Return(existingProduct, nil)
	mockProductRepo.On("UpdateProductStock", ctx, (*gorm.DB)(nil), productID.String(), 50).
		Return(nil)

	// Execute
	result, err := productService.UpdateStock(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 150, result.Stock)
	mockProductRepo.AssertExpectations(t)
}

func TestUpdateStock_InsufficientStock(t *testing.T) {
	// Setup
	mockProductRepo := new(mockProductRepository)
	mockCategoryRepo := new(mockCategoryRepository)
	mockProductQ := new(mockProductQuery)
	mockCategoryQ := new(mockCategoryQuery)
	mockTxRepo := new(mockTxRepository)

	productService := service.NewProductService(
		mockProductRepo, mockCategoryRepo, mockProductQ, mockCategoryQ, mockTxRepo,
	)

	ctx := context.Background()
	productID := uuid.New()

	existingProduct := entity.Product{
		ID:    productID,
		Name:  "Test Product",
		Stock: 10,
	}

	req := dto.ProductStockUpdateRequest{
		ID:       productID.String(),
		Quantity: -20, // Trying to reduce more than available
	}

	// Expectations
	mockProductRepo.On("GetProductByID", ctx, (*gorm.DB)(nil), productID.String(), []string(nil)).
		Return(existingProduct, nil)

	// Execute
	_, err := productService.UpdateStock(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errs.ErrInsufficientStock, err)
	mockProductRepo.AssertExpectations(t)
}

func TestGetLowStockProducts_Success(t *testing.T) {
	// Setup
	mockProductRepo := new(mockProductRepository)
	mockCategoryRepo := new(mockCategoryRepository)
	mockProductQ := new(mockProductQuery)
	mockCategoryQ := new(mockCategoryQuery)
	mockTxRepo := new(mockTxRepository)

	productService := service.NewProductService(
		mockProductRepo, mockCategoryRepo, mockProductQ, mockCategoryQ, mockTxRepo,
	)

	ctx := context.Background()
	threshold := 10

	lowStockProducts := []entity.Product{
		{ID: uuid.New(), Name: "Product 1", Stock: 5},
		{ID: uuid.New(), Name: "Product 2", Stock: 3},
	}

	// Expectations
	mockProductQ.On("GetLowStockProducts", ctx, threshold).
		Return(lowStockProducts, nil)

	// Execute
	result, err := productService.GetLowStockProducts(ctx, threshold)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockProductQ.AssertExpectations(t)
}
