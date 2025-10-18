package queryiface

import (
	"context"

	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"github.com/zetsux/gin-gorm-clean-starter/support/base"
)

type UserQuery interface {
	GetAllUsers(ctx context.Context, req base.GetsRequest) ([]entity.User, int64, int64, error)
}
