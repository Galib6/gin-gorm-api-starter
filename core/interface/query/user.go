package queryiface

import (
	"context"

	"github.com/zetsux/gin-gorm-api-starter/core/entity"
	"github.com/zetsux/gin-gorm-api-starter/core/helper/dto"
	"github.com/zetsux/gin-gorm-api-starter/support/base"
)

type UserQuery interface {
	GetAllUsers(ctx context.Context, req dto.UserGetsRequest) ([]entity.User, base.PaginationResponse, error)
}
