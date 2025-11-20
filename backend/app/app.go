package app

import (
	"database/sql"
	"structured-notes/repositories"
	"structured-notes/services"
	"structured-notes/utils"
)

type Config struct {
	Port     int
	Database struct {
		Host   string
		Port   int
		Name   string
		Driver string
	}
	Media struct {
		MaxSize              float64
		MaxUploadsSize       float64
		SupportedTypesImages []string
		SupportedTypes       []string
	}
	Auth struct {
		AccessTokenExpiry  int
		RefreshTokenExpiry int
	}
}

type App struct {
	DB        *sql.DB
	Snowflake *utils.Snowflake
	Config    Config
	Services  *services.ServiceManager
	Repos     *repositories.RepositoryManager
}

func InitApp(config Config) *App {
	var app App

	return &app
}
