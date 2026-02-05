package migrations

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// GetDB returns a database connection for migrations
func GetDB() (*sql.DB, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Println("Warning: No .env file loaded")
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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations runs all pending migrations
func RunMigrations(direction string) error {
	db, err := GetDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// Set the migrations directory
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	migrationsDir := "database/migrations"

	switch direction {
	case "up":
		if err := goose.Up(db, migrationsDir); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	case "down":
		if err := goose.Down(db, migrationsDir); err != nil {
			return fmt.Errorf("failed to rollback migration: %w", err)
		}
	case "status":
		if err := goose.Status(db, migrationsDir); err != nil {
			return fmt.Errorf("failed to get migration status: %w", err)
		}
	case "reset":
		if err := goose.Reset(db, migrationsDir); err != nil {
			return fmt.Errorf("failed to reset migrations: %w", err)
		}
	default:
		return fmt.Errorf("unknown migration direction: %s", direction)
	}

	return nil
}

// CreateMigration creates a new migration file
func CreateMigration(name string, migrationType string) error {
	migrationsDir := "database/migrations"

	if migrationType == "" {
		migrationType = "sql"
	}

	if err := goose.Create(nil, migrationsDir, name, migrationType); err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	return nil
}
