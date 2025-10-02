package factory

import (
	"testing"

	"github.com/zetsux/gin-gorm-clean-starter/core/repository"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"gorm.io/gorm"
)

func NewUserRepository(t *testing.T, db *gorm.DB) repository.UserRepository {
	t.Helper()
	txr := repository.NewTxRepository(db)
	return repository.NewUserRepository(txr)
}

func NewUserService(t *testing.T, db *gorm.DB) service.UserService {
	t.Helper()
	ur := NewUserRepository(t, db)
	return service.NewUserService(ur)
}
