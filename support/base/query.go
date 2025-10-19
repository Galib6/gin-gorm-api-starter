package base

import (
	"math"

	errs "github.com/zetsux/gin-gorm-clean-starter/core/helper/errors"
	"gorm.io/gorm"
)

func GetWithPagination[T any](stmt *gorm.DB, req PaginationRequest,
) (data []T, paginationResp PaginationResponse, err error) {
	var totalCount int64
	if err := stmt.Count(&totalCount).Error; err != nil {
		return nil, paginationResp, err
	}

	if req.Sort != "" {
		sortBy := req.Sort
		if sortBy[0] == '-' {
			sortBy = sortBy[1:] + " DESC"
		}
		stmt = stmt.Order(sortBy)
	}

	if req.PerPage == 0 {
		err = stmt.Find(&data).Error
		return data, paginationResp, err
	}

	lastPage := int64(math.Ceil(float64(totalCount) / float64(req.PerPage)))
	if req.Page <= 0 || int64(req.Page) > lastPage {
		return nil, paginationResp, errs.ErrInvalidPage
	}

	offset := (req.Page - 1) * req.PerPage
	err = stmt.Offset(offset).Limit(req.PerPage).Find(&data).Error

	paginationResp.Page = int64(req.Page)
	paginationResp.PerPage = int64(req.PerPage)
	paginationResp.LastPage = lastPage
	paginationResp.Total = totalCount

	return data, paginationResp, err
}
