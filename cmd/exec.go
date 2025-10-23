package cmd

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

func Execute(db *gorm.DB) {
	if len(os.Args) < 2 {
		return
	}

	command := os.Args[1]
	switch command {
	case "setup":
		RunMigration(db)
		RunSeeder(db)
		os.Exit(0)
	case "migrate":
		RunMigration(db)
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

func printHelp() {
	fmt.Println(`
📘 Available commands:
	setup     → Run migration + seed
	migrate   → Run migration only
	seed      → Run seeder only
	`)
}
