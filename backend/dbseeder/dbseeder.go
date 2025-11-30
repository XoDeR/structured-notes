package dbseeder

import (
	"structured-notes/logger"
	"structured-notes/services"
)

func Seed(services *services.ServiceManager) error {
	// insert test data
	logger.Info("Db Seeder: inserting test data")
	return nil
}

func Truncate(services *services.ServiceManager) error {
	// truncate tables
	logger.Info("Db Seeder: Truncating tables")
	return nil
}
