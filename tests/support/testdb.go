package support

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/zetsux/gin-gorm-clean-starter/core/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewTestDB starts an embedded postgres, migrates schemas, and returns a *gorm.DB.
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	port := uint32(5432 + rand.Intn(1000))
	base := t.TempDir()
	runtimePath := filepath.Join(base, "rt")
	dataPath := filepath.Join(base, "data")
	_ = os.MkdirAll(runtimePath, 0o755)
	_ = os.MkdirAll(dataPath, 0o755)
	cfg := embeddedpostgres.DefaultConfig().
		Port(port).
		Database("testdb").
		Username("test").
		Password("test").
		RuntimePath(runtimePath).
		DataPath(dataPath).
		Version(embeddedpostgres.V15)
	pg := embeddedpostgres.NewDatabase(cfg)
	if err := pg.Start(); err != nil {
		t.Fatalf("failed to start embedded postgres: %v", err)
	}
	// ensure stop
	t.Cleanup(func() { _ = pg.Stop() })

	dsn := fmt.Sprintf("host=localhost user=test password=test dbname=testdb"+
		" port=%d sslmode=disable TimeZone=Asia/Jakarta", port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open postgres: %v", err)
	}
	// enable uuid extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		t.Fatalf("failed to enable uuid extension: %v", err)
	}
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}
	return db
}
