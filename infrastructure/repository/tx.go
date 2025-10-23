package repository

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type txRepository struct {
	db *gorm.DB
}

func NewTxRepository(db *gorm.DB) *txRepository {
	return &txRepository{db: db}
}

func (rp txRepository) DB() *gorm.DB {
	return rp.db
}

func (rp txRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := rp.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (rp txRepository) CommitOrRollbackTx(ctx context.Context, tx *gorm.DB, err error) {
	if err != nil {
		log.Println("Error occurred: ", err)
		tx.WithContext(ctx).Debug().Rollback()
		return
	}

	err = tx.WithContext(ctx).Commit().Error
	if err != nil {
		log.Println("Commit failed: ", err)
		return
	}
	log.Println("Committed successfully")
}
