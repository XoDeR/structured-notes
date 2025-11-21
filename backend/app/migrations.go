package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(config *Config) {
	workingDir, _ := os.Getwd()
	absPath := filepath.Join(workingDir, os.Getenv("CONFIG_CPWD"), "migrations")
	absPath = filepath.ToSlash(absPath)
	db := DBConnection(*config, true) // multiStatements are used for migrations
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create mysql driver: %v", err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", absPath),
		"structured-notes",
		driver,
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create migrate instance: %v", err))
	}
	defer m.Close()

	if err := m.Up(); err != nil && err.Error() != "no change" {
		panic(fmt.Sprintf("Failed to apply migrations: %v", err))
	}
}
