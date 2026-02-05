package database

import (
	"fmt"

	"myapp/core/entity"
	"myapp/database/seeder"

	"gorm.io/gorm"
)

func DBMigrate(db *gorm.DB) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		fmt.Println("Failed to create uuid-ossp extension:", err)
		panic(err)
	}

	err := db.AutoMigrate(
		entity.User{},
		entity.Category{},
		entity.Product{},
	)

	if err != nil {
		fmt.Println("Failed to migrate database: ", err)
		panic(err)
	}
}

func DBSeed(db *gorm.DB) {
	if err := seeder.UserSeeder(db); err != nil {
		fmt.Println("Failed to seed users: ", err)
		panic(err)
	}

	if err := seeder.ProductSeeder(db); err != nil {
		fmt.Println("Failed to seed products: ", err)
		panic(err)
	}
}
