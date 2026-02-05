package query

import (
	"fmt"
	"math"
	"strings"

	errs "myapp/core/helper/errors"
	"myapp/support/base"

	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func applySorting(stmt *gorm.DB, allowedSorts []string, sort string) (*gorm.DB, error) {
	col := sort
	direction := " ASC"

	if sort == "" {
		col = allowedSorts[0]
		return stmt.Order(col + direction), nil
	}

	if strings.HasPrefix(sort, "-") {
		col = sort[1:]
		direction = " DESC"
	}
	if !slices.Contains(allowedSorts, col) {
		return nil, fmt.Errorf("%w: column '%s' (allowed values: %s)",
			errs.ErrInvalidSort, col, strings.Join(allowedSorts, ", "))
	}

	return stmt.Order(col + direction), nil
}

func applyIncludes(stmt *gorm.DB, allowedIncludes []string, includes string) (*gorm.DB, error) {
	allowedValues := "none"
	if len(allowedIncludes) > 0 {
		allowedValues = strings.Join(allowedIncludes, ", ")
	}

	for _, include := range strings.Split(includes, ",") {
		include = strings.TrimSpace(include)
		if include == "" {
			continue
		}

		if !slices.Contains(allowedIncludes, include) {
			return nil, fmt.Errorf("%w: column '%s' (allowed values: %s)",
				errs.ErrInvalidInclude, include, allowedValues)
		}
		stmt = stmt.Preload(include)
	}
	return stmt, nil
}

func GetWithPagination[T any](
	stmt *gorm.DB, req base.PaginationRequest,
	allowedSorts []string, allowedIncludes []string,
) (data []T, paginationResp base.PaginationResponse, err error) {
	stmt, err = applySorting(stmt, allowedSorts, req.Sort)
	if err != nil {
		return nil, paginationResp, err
	}

	stmt, err = applyIncludes(stmt, allowedIncludes, req.Includes)
	if err != nil {
		return nil, paginationResp, err
	}

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
