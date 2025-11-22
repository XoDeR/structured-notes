package repositories

import (
	"database/sql"
	"fmt"
	"structured-notes/logger"
	"sync"
)

type RepositoryManager struct {
	db          *sql.DB
	statements  map[string]*sql.Stmt
	stmtMutex   sync.RWMutex
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

func (rm *RepositoryManager) PrepareStatement(key string, query string) (*sql.Stmt, error) {
	rm.stmtMutex.Lock()
	defer rm.stmtMutex.Unlock()

	// Check if statement already exists
	if stmt, exists := rm.statements[key]; exists {
		return stmt, nil
	}

	// Prepare new statement
	stmt, err := rm.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement '%s': %w", key, err)
	}

	// Cache the statement
	rm.statements[key] = stmt
	return stmt, nil
}

func (rm *RepositoryManager) GetStatement(key string) (*sql.Stmt, error) {
	rm.stmtMutex.RLock()
	defer rm.stmtMutex.RUnlock()

	stmt, exists := rm.statements[key]
	if !exists {
		return nil, fmt.Errorf("Prepared statement '%s' not found", key)
	}
	return stmt, nil
}
