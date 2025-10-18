package repositoryiface

import (
	"context"

	"gorm.io/gorm"
)

type TxRepository interface {
	// db
	DB() *gorm.DB

	// tx
	BeginTx(ctx context.Context) (*gorm.DB, error)
	CommitOrRollbackTx(ctx context.Context, tx *gorm.DB, err error)
}
