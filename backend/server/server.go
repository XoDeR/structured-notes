package server

import (
	"fmt"
	"os"
	"path/filepath"
	"structured-notes/app"
	"structured-notes/logger"
	"structured-notes/router"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

func SetupServer() (*gin.Engine, *app.App) {
	workingDir, err := os.Getwd()
	if err != nil {
		logger.Error("Error getting cwd: " + err.Error())
		os.Exit(1)
	}
	absPath := filepath.Join(workingDir, fmt.Sprintf("%sconfig.toml", os.Getenv("CONFIG_CPWD")))
	config := app.Config{}
	_, err = toml.DecodeFile(absPath, &config)
	if err != nil {
		logger.Error("Error loading config: " + err.Error())
		os.Exit(1)
	}
	logger.Success("Loaded configuration from: " + absPath + " successfully")

	application := app.InitApp(config)

	appRouter := router.InitRouter(application)

	return appRouter, application
}
