package app

import (
	"database/sql"
	"fmt"
	"os"
	"structured-notes/logger"
	"time"
)

func DBConnection(config Config, multiStatements bool) (connection *sql.DB) {
	Driver := config.Database.Driver
	User := os.Getenv("DATABASE_USER")
	Password := os.Getenv("DATABASE_PASSWORD")

	// If no env vars fallback to values from config file
	Host := os.Getenv("DATABASE_HOST")
	if Host == "" {
		Host = config.Database.Host
	}

	Port := os.Getenv("DATABASE_PORT")
	if Port == "" {
		Port = fmt.Sprint(config.Database.Port)
	}

	Database := os.Getenv("DATABASE_NAME")
	if Database == "" {
		Database = config.Database.Name
	}

	multiStatementsConfig := "?multiStatements=false"
	if multiStatements {
		multiStatementsConfig = "?multiStatements=true"
	}

	userPassword := User
	if len(Password) > 0 {
		userPassword += ":" + Password
	}

	connection, err := sql.Open(Driver, userPassword+"@tcp("+Host+":"+Port+")/"+Database+multiStatementsConfig)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %v", err))
	}
	logger.Info("Connected to db")
	connection.SetConnMaxLifetime(time.Minute * 3)
	connection.SetMaxOpenConns(10)
	connection.SetMaxIdleConns(10)
	return connection
}
