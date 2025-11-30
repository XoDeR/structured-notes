package main

import (
	"flag"
	"fmt"
	"os"
	"structured-notes/dbseeder"
	"structured-notes/logger"
	"structured-notes/server"

	"github.com/joho/godotenv"
)

var (
	devMode    = flag.Bool("dev", false, "Run in development mode")
	dbSeed     = flag.Bool("db-seed", false, "Seed the database with test data (dev mode only)")
	dbTruncate = flag.Bool("db-truncate", false, "Truncate the database tables (dev mode only)")
)

func main() {
	logger.Info("Initializing SN backend...")

	// get env vars
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
		os.Exit(1)
	}

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		logger.Error("BACKEND_PORT environment variable not set")
		os.Exit(1)
	}
	appRouter, application := server.SetupServer()

	// Seed/truncate in dev mode only
	flag.Parse()

	// Only allow seeding/truncating in dev mode
	if *devMode {
		if *dbTruncate {
			if err := dbseeder.Truncate(application.Services); err != nil {
				fmt.Println("truncate failed:", err)
				os.Exit(1)
			}
			fmt.Println("Database truncated successfully")
		}

		if *dbSeed {
			if err := dbseeder.Seed(application.Services); err != nil {
				fmt.Println("seed failed: ", err)
				os.Exit(1)
			}
			fmt.Println("Database seeded successfully")
		}
	}
	//-- Seed/truncate

	logger.Info("Starting server on port: " + port)
	defer application.DB.Close()

	appRouter.Run(":" + port)
}
