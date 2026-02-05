package query

import (
	"context"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	"myapp/support/base"

	"gorm.io/gorm"
)

var userAllowedSorts = []string{"id", "name", "email", "created_at", "updated_at"}
var userAllowedIncludes = []string{}

type userQuery struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) *userQuery {
	return &userQuery{db: db}
}

func (qr *userQuery) GetAllUsers(ctx context.Context, req dto.UserGetsRequest,
) ([]entity.User, base.PaginationResponse, error) {
	stmt := qr.db.WithContext(ctx).Debug().Model(&entity.User{})

	if req.ID != "" {
		stmt = stmt.Where("id = ?", req.ID)
	}

	if req.Role != "" {
		stmt = stmt.Where("role = ?", req.Role)
	}

	if req.Search != "" {
		search := "%" + req.Search + "%"
		stmt = stmt.Where("name ILIKE ? OR email ILIKE ?", search, search)
	}

	users, pageResp, err := GetWithPagination[entity.User](stmt,
		req.PaginationRequest, userAllowedSorts, userAllowedIncludes)
	if err != nil {
		return nil, pageResp, err
	}
	return users, pageResp, nil
}
