package database

import (
	"fmt"

	"github.com/zetsux/gin-gorm-api-starter/core/entity"
	"github.com/zetsux/gin-gorm-api-starter/database/seeder"
	"gorm.io/gorm"
)

func DBMigrate(db *gorm.DB) {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		fmt.Println("Failed to create uuid-ossp extension:", err)
		panic(err)
	}

	err := db.AutoMigrate(
		entity.User{},
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
}
