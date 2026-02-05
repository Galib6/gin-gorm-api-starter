package query

import (
	"context"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	"myapp/support/base"

	"gorm.io/gorm"
)

var categoryAllowedSorts = []string{"id", "name", "created_at", "updated_at"}
var categoryAllowedIncludes = []string{"Products"}

type categoryQuery struct {
	db *gorm.DB
}

func NewCategoryQuery(db *gorm.DB) *categoryQuery {
	return &categoryQuery{db: db}
}

func (qr *categoryQuery) GetAllCategories(ctx context.Context, req dto.CategoryGetsRequest,
) ([]entity.Category, base.PaginationResponse, error) {
	stmt := qr.db.WithContext(ctx).Debug().Model(&entity.Category{})

	if req.ID != "" {
		stmt = stmt.Where("id = ?", req.ID)
	}

	if req.Search != "" {
		search := "%" + req.Search + "%"
		stmt = stmt.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	categories, pageResp, err := GetWithPagination[entity.Category](stmt,
		req.PaginationRequest, categoryAllowedSorts, categoryAllowedIncludes)
	if err != nil {
		return nil, pageResp, err
	}
	return categories, pageResp, nil
}
