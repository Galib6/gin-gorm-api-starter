package query

import (
	"context"
	"errors"
	"math"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func NewUserQuery(db *gorm.DB) *userQuery {
	return &userQuery{db: db}
}

func (uq *userQuery) GetAllUsers(ctx context.Context, req base.GetsRequest) ([]entity.User, int64, int64, error) {
	var err error
	var users []entity.User
	var total int64

	stmt := uq.db.WithContext(ctx).Debug()
	if req.Search != "" {
		searchQuery := "%" + req.Search + "%"
		err = uq.db.WithContext(ctx).Model(&entity.User{}).
			Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery).
			Count(&total).Error

		if err != nil {
			return nil, 0, 0, err
		}
		stmt = stmt.Where("name ILIKE ? OR email ILIKE ?", searchQuery, searchQuery)
	} else {
		err = uq.db.WithContext(ctx).Model(&entity.User{}).Count(&total).Error
		if err != nil {
			return nil, 0, 0, err
		}
	}

	if req.Sort != "" {
		stmt = stmt.Order(req.Sort)
	}

	lastPage := int64(math.Ceil(float64(total) / float64(req.PerPage)))
	if req.PerPage == 0 {
		err = stmt.Find(&users).Error
	} else {
		if req.Page <= 0 || int64(req.Page) > lastPage {
			return nil, 0, 0, errs.ErrInvalidPage
		}
		err = stmt.Offset(((req.Page - 1) * req.PerPage)).Limit(req.PerPage).Find(&users).Error
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		return users, 0, 0, err
	}
	return users, lastPage, total, nil
}
