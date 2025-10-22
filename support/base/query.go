package base

import (
	"math"
	"strings"

	errs "github.com/zetsux/gin-gorm-api-starter/core/helper/errors"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func applySorting(stmt *gorm.DB, allowedSorts []string, sort string) *gorm.DB {
	col := sort
	direction := " ASC"

	if strings.HasPrefix(sort, "-") {
		col = sort[1:]
		direction = " DESC"
	}

	if !slices.Contains(allowedSorts, col) {
		col = allowedSorts[0]
		direction = " ASC"
	}

	stmt = stmt.Order(col + direction)
	return stmt
}

func GetWithPagination[T any](stmt *gorm.DB, req PaginationRequest, allowedSorts []string,
) (data []T, paginationResp PaginationResponse, err error) {
	stmt = applySorting(stmt, allowedSorts, req.Sort)

	if req.PerPage == 0 {
		err = stmt.Find(&data).Error
		return data, paginationResp, err
	}

	var totalCount int64
	if err := stmt.Count(&totalCount).Error; err != nil {
		return nil, paginationResp, err
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
