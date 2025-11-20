package main

import (
	"os"
	"structured-notes/logger"
	"structured-notes/server"

	"github.com/joho/godotenv"
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

	logger.Info("Starting server on port: " + port)
	defer application.DB.Close()

	appRouter.Run(":" + port)
}
