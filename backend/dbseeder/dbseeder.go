package dbseeder

import (
	"structured-notes/logger"
	"structured-notes/services"
)

func Seed(services *services.ServiceManager) error {
	// insert test data
	logger.Info("Db Seeder: inserting test data")

	/*
		// Create admin
		{
			username := "admin"
			password := "test"
			services.User.CreateUser(username, "Admin", "Admin", "", "admin@example.com", password)
		}

		// Create user
		{
			username := "johnj"
			password := "test"
			services.User.CreateUser(username, "John", "Doe", "", "john@example.com", password)
		}

		// Create other 10 users
		{
			count := 10
			for i := range count {
				username := "bill" + strconv.Itoa(i)
				password := "test"
				services.User.CreateUser(username, "Bill"+strconv.Itoa(i), "Doe"+strconv.Itoa(i), "", "bill"+strconv.Itoa(i)+"@example.com", password)
			}
		}

		// Create 10 categories by user johnj
		// Create 10 documents by user johnj
		// Create 10 media files by user johnj
	*/

	return nil
}

func Truncate(services *services.ServiceManager) error {
	// truncate tables
	logger.Info("Db Seeder: Truncating tables")

	// connection_logs

	// permissions

	// nodes

	// sessions

	// users

	return nil
}
