package repositories

import (
	"database/sql"
	"structured-notes/logger"
)

type RepositoryManager struct {
	db          *sql.DB
	initialized bool
}

func NewRepositoryManager(db *sql.DB) (*RepositoryManager, error) {
	rm := &RepositoryManager{
		db: db,
	}

	rm.initialized = true
	logger.Success("Repository manager init success")
	return rm, nil
}
