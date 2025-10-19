package query

import (
	"context"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) *userQuery {
	return &userQuery{db: db}
}

func (uq *userQuery) GetAllUsers(ctx context.Context, req dto.UserGetsRequest,
) ([]entity.User, base.PaginationResponse, error) {
	stmt := uq.db.WithContext(ctx).Debug().Model(&entity.User{})

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

	users, pageResp, err := base.GetWithPagination[entity.User](stmt, req.PaginationRequest)
	if err != nil {
		return nil, pageResp, err
	}
	return users, pageResp, nil
}
