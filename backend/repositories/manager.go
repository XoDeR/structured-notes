package repositories

import (
	"database/sql"
	"fmt"
	"structured-notes/logger"
	"sync"
)

type RepositoryManager struct {
	db          *sql.DB
	User        UserRepository
	Node        NodeRepository
	Session     SessionRepository
	Permission  PermissionRepository
	Log         LogRepository
	statements  map[string]*sql.Stmt
	stmtMutex   sync.RWMutex
	initialized bool
}

func NewRepositoryManager(db *sql.DB) (*RepositoryManager, error) {
	rm := &RepositoryManager{
		db:         db,
		statements: make(map[string]*sql.Stmt),
	}

	if err := rm.initializeRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	rm.initialized = true
	logger.Success("Repository manager init success")
	return rm, nil
}

func (rm *RepositoryManager) initializeRepositories() error {
	var err error

	rm.User, err = NewUserRepository(rm.db, rm)
	if err != nil {
		return fmt.Errorf("failed to initialize user repository: %w", err)
	}

	rm.Node, err = NewNodeRepository(rm.db, rm)
	if err != nil {
		return fmt.Errorf("failed to initialize node repository: %w", err)
	}

	rm.Session, err = NewSessionRepository(rm.db, rm)
	if err != nil {
		return fmt.Errorf("failed to initialize session repository: %w", err)
	}

	rm.Permission, err = NewPermissionRepository(rm.db, rm)
	if err != nil {
		return fmt.Errorf("failed to initialize permission repository: %w", err)
	}

	rm.Log, err = NewLogRepository(rm.db, rm)
	if err != nil {
		return fmt.Errorf("failed to initialize log repository: %w", err)
	}

	return nil
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
		return nil, fmt.Errorf("prepared statement '%s' not found", key)
	}
	return stmt, nil
}

func (rm *RepositoryManager) Close() error {
	rm.stmtMutex.Lock()
	defer rm.stmtMutex.Unlock()

	var errors []error
	for key, stmt := range rm.statements {
		if err := stmt.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close statement '%s': %w", key, err))
		}
	}

	rm.statements = make(map[string]*sql.Stmt)
	rm.initialized = false

	if len(errors) > 0 {
		return fmt.Errorf("errors occurred while closing statements: %v", errors)
	}

	logger.Success("Repository manager closed successfully")
	return nil
}

func (rm *RepositoryManager) IsInitialized() bool {
	return rm.initialized
}

func (rm *RepositoryManager) GetStatementCount() int {
	rm.stmtMutex.RLock()
	defer rm.stmtMutex.RUnlock()
	return len(rm.statements)
}
