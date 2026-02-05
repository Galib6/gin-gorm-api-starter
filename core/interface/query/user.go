package queryiface

import (
	"context"

	"myapp/core/entity"
	"myapp/core/helper/dto"
	"myapp/support/base"
)

type UserQuery interface {
	GetAllUsers(ctx context.Context, req dto.UserGetsRequest) ([]entity.User, base.PaginationResponse, error)
}
