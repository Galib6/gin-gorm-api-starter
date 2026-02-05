package database

import (
	"fmt"

	"myapp/database/seeder"

	"gorm.io/gorm"
)

// DBSeed runs all database seeders
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
