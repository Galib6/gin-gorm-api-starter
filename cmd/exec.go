package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"myapp/database/migrations"

	"gorm.io/gorm"
)

func Execute(db *gorm.DB) {
	if len(os.Args) < 2 {
		return
	}

	command := os.Args[1]
	switch command {
	case "setup":
		RunGooseMigration("up")
		RunSeeder(db)
		os.Exit(0)
	case "migrate":
		RunGooseMigration("up")
		os.Exit(0)
	case "migrate:generate":
		if len(os.Args) < 3 {
			fmt.Println("❌ Usage: go run main.go migrate:generate <migration_name>")
			os.Exit(1)
		}
		name := os.Args[2]
		cmd := exec.Command("atlas", "migrate", "diff", name, "--env", "local")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("❌ Atlas failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Migration generated successfully!")
		os.Exit(0)
	case "migrate:down":
		RunGooseMigration("down")
		os.Exit(0)
	case "migrate:status":
		RunGooseMigration("status")
		os.Exit(0)
	case "migrate:reset":
		RunGooseMigration("reset")
		os.Exit(0)
	case "migrate:create":
		if len(os.Args) < 3 {
			fmt.Println("❌ Please provide migration name: go run main.go migrate:create <name>")
			os.Exit(1)
		}
		migType := "sql"
		if len(os.Args) >= 4 {
			migType = os.Args[3]
		}
		if err := migrations.CreateMigration(os.Args[2], migType); err != nil {
			fmt.Printf("❌ Failed to create migration: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Migration created successfully")
		os.Exit(0)
	case "seed":
		RunSeeder(db)
		os.Exit(0)
	default:
		fmt.Printf("⚠️  Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func RunGooseMigration(direction string) {
	fmt.Printf("🔄 Running migration: %s\n", direction)
	if err := migrations.RunMigrations(direction); err != nil {
		fmt.Printf("❌ Migration failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Migration completed successfully")
}

func printHelp() {
	fmt.Println(`
📘 Available commands:
	setup              → Run Migrations + Seed
	migrate            → Run all pending migrations (Goose up)
	migrate:down       → Rollback last migration
	migrate:status     → Show migration status
	migrate:reset      → Reset all migrations
	migrate:generate   → ✨ Auto-Generate SQL from GORM Models (Atlas)
	migrate:create     → Create empty migration file (manual SQL)
	seed               → Run seeder only
	`)
}
