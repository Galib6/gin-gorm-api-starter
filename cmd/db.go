package cmd

import (
	"fmt"

	"myapp/database"

	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB) {
	fmt.Println("🌱 Starting database seeding...")
	database.DBSeed(db)
	fmt.Println("✅ Database seeding completed successfully!")
}
