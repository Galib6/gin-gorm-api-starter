package cmd

import (
	"fmt"

	"github.com/zetsux/gin-gorm-api-starter/database"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	fmt.Println("🚀 Starting database migration...")
	database.DBMigrate(db)
	fmt.Println("✅ Database migration completed successfully!")
}

func RunSeeder(db *gorm.DB) {
	fmt.Println("🌱 Starting database seeding...")
	database.DBSeed(db)
	fmt.Println("✅ Database seeding completed successfully!")
}
