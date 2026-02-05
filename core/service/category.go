package service

import (
	"context"
	"reflect"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	errs "myapp/core/helper/errors"
	queryiface "myapp/core/interface/query"
	repositoryiface "myapp/core/interface/repository"
	"myapp/support/base"
	"myapp/support/constant"
)

type categoryService struct {
	categoryRepository repositoryiface.CategoryRepository
	productRepository  repositoryiface.ProductRepository
	categoryQuery      queryiface.CategoryQuery
	productQuery       queryiface.ProductQuery
}

type CategoryService interface {
	CreateCategory(ctx context.Context, req dto.CategoryCreateRequest) (dto.CategoryResponse, error)
	GetAllCategories(ctx context.Context, req dto.CategoryGetsRequest) ([]dto.CategoryResponse, base.PaginationResponse, error)
	GetCategoryByID(ctx context.Context, id string) (dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, req dto.CategoryUpdateRequest) (dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) error
}

func NewCategoryService(
	categoryR repositoryiface.CategoryRepository,
	productR repositoryiface.ProductRepository,
	categoryQ queryiface.CategoryQuery,
	productQ queryiface.ProductQuery,
) CategoryService {
	return &categoryService{
		categoryRepository: categoryR,
		productRepository:  productR,
		categoryQuery:      categoryQ,
		productQuery:       productQ,
	}
}

func (sv *categoryService) toCategoryResponse(category entity.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
	}
}

func (sv *categoryService) CreateCategory(ctx context.Context, req dto.CategoryCreateRequest) (dto.CategoryResponse, error) {
	// Check if category name already exists
	existing, err := sv.categoryRepository.GetCategoryByPrimaryKey(ctx, nil, constant.DBAttrName, req.Name)
	if err != nil && err != errs.ErrCategoryNotFound {
		return dto.CategoryResponse{}, err
	}
	if !reflect.DeepEqual(existing, entity.Category{}) {
		return dto.CategoryResponse{}, errs.ErrCategoryNameExists
	}

	category := entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	newCategory, err := sv.categoryRepository.CreateCategory(ctx, nil, category)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	return sv.toCategoryResponse(newCategory), nil
}

func (sv *categoryService) GetAllCategories(ctx context.Context, req dto.CategoryGetsRequest) (
	categoriesResp []dto.CategoryResponse, pageResp base.PaginationResponse, err error) {

	categories, pageResp, err := sv.categoryQuery.GetAllCategories(ctx, req)
	if err != nil {
		return []dto.CategoryResponse{}, base.PaginationResponse{}, err
	}

	for _, category := range categories {
		categoriesResp = append(categoriesResp, sv.toCategoryResponse(category))
	}
	return categoriesResp, pageResp, nil
}

func (sv *categoryService) GetCategoryByID(ctx context.Context, id string) (dto.CategoryResponse, error) {
	category, err := sv.categoryRepository.GetCategoryByID(ctx, nil, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	return sv.toCategoryResponse(category), nil
}

func (sv *categoryService) UpdateCategory(ctx context.Context, req dto.CategoryUpdateRequest) (dto.CategoryResponse, error) {
	category, err := sv.categoryRepository.GetCategoryByID(ctx, nil, req.ID)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	// Check if new name already exists (if changing)
	if req.Name != "" && req.Name != category.Name {
		existing, err := sv.categoryRepository.GetCategoryByPrimaryKey(ctx, nil, constant.DBAttrName, req.Name)
		if err != nil && err != errs.ErrCategoryNotFound {
			return dto.CategoryResponse{}, err
		}
		if !reflect.DeepEqual(existing, entity.Category{}) {
			return dto.CategoryResponse{}, errs.ErrCategoryNameExists
		}
	}

	categoryEdit := entity.Category{
		ID:          category.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	err = sv.categoryRepository.UpdateCategory(ctx, nil, categoryEdit)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	return sv.toCategoryResponse(categoryEdit), nil
}

func (sv *categoryService) DeleteCategory(ctx context.Context, id string) error {
	_, err := sv.categoryRepository.GetCategoryByID(ctx, nil, id)
	if err != nil {
		return err
	}

	// Check if category has products
	products, _, err := sv.productQuery.GetAllProducts(ctx, dto.ProductGetsRequest{
		CategoryID: id,
		PaginationRequest: base.PaginationRequest{
			Page:    1,
			PerPage: 1,
		},
	})
	if err != nil {
		return err
	}
	if len(products) > 0 {
		return errs.ErrCategoryHasProducts
	}

	return sv.categoryRepository.DeleteCategoryByID(ctx, nil, id)
}
