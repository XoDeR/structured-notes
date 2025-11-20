package main

import (
	"fmt"
	"os"
	"structured-notes/logger"

	"github.com/joho/godotenv"
)

func main() {
	logger.Info("Initializing SN backend...")

	// get env vars
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	testEnvVar := os.Getenv("TEST")
	fmt.Println("TEST:", testEnvVar)

	fmt.Println("Structured Notes backend is running...")
}
