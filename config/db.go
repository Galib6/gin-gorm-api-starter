package config

import (
	"fmt"
	"os"

	migration "github.com/zetsux/gin-gorm-clean-starter/database"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBSetup() *gorm.DB {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	prefix := ""
	if os.Getenv("USE_DOCKER") == "true" {
		prefix = "DOCKER_"
	}

	dbUser := os.Getenv(prefix + "DB_USER")
	dbPass := os.Getenv(prefix + "DB_PASS")
	dbHost := os.Getenv(prefix + "DB_HOST")
	dbName := os.Getenv(prefix + "DB_NAME")
	dbPort := os.Getenv(prefix + "DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	migration.DBMigrate(db)
	return db
}

func DBClose(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	dbSQL.Close()
}
