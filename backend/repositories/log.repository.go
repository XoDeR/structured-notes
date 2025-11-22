package repositories

import (
	"database/sql"
	"fmt"
	"structured-notes/models"
	"structured-notes/types"
	"time"
)

type LogRepository interface {
	GetByUserID(userId types.Snowflake) ([]*models.Log, error)
	GetLastByUserID(userId types.Snowflake) (*models.Log, error)
	Create(log *models.Log) error
	DeleteOld() error
}

type LogRepositoryImpl struct {
	db      *sql.DB
	manager *RepositoryManager
}

const (
	stmtLogGetByUserID       = "log_get_by_user_id"
	stmtLogGetLastConnection = "log_get_last_connection"
	stmtLogCreateConnection  = "log_create_connection"
	stmtLogDeleteOld         = "log_delete_old"
)

func NewLogRepository(db *sql.DB, manager *RepositoryManager) (LogRepository, error) {
	repo := &LogRepositoryImpl{
		db:      db,
		manager: manager,
	}

	if err := repo.prepareStatements(); err != nil {
		return nil, fmt.Errorf("failed to prepare log statements: %w", err)
	}

	return repo, nil
}

func (r *LogRepositoryImpl) prepareStatements() error {
	statements := map[string]string{
		stmtLogGetByUserID: `
			SELECT id, user_id, ip_adress, timestamp, type, location, user_agent
			FROM connections_logs
			WHERE user_id = ?
			ORDER BY timestamp DESC`,

		stmtLogGetLastConnection: `
			SELECT id, user_id, ip_adress, timestamp, type, location, user_agent
			FROM connections_logs
			WHERE user_id = ? AND type = 'login'
			ORDER BY timestamp DESC
			LIMIT 1`,

		stmtLogCreateConnection: `
			INSERT INTO connections_logs (id, user_id, ip_adress, timestamp, type, location, user_agent)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,

		stmtLogDeleteOld: `
			DELETE FROM connections_logs
			WHERE timestamp < ?`,
	}

	// Prepare all statements
	for key, query := range statements {
		if _, err := r.manager.PrepareStatement(key, query); err != nil {
			return err
		}
	}

	return nil
}

func (r *LogRepositoryImpl) GetByUserID(userId types.Snowflake) ([]*models.Log, error) {
	stmt, err := r.manager.GetStatement(stmtLogGetByUserID)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	logs := make([]*models.Log, 0)
	for rows.Next() {
		var log models.Log
		err := rows.Scan(
			&log.Id,
			&log.UserId,
			&log.IpAddr,
			&log.Timestamp,
			&log.Type,
			&log.Location,
			&log.UserAgent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, &log)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, nil
}

func (r *LogRepositoryImpl) GetLastByUserID(userId types.Snowflake) (*models.Log, error) {
	stmt, err := r.manager.GetStatement(stmtLogGetLastConnection)
	if err != nil {
		return nil, err
	}

	var log models.Log
	err = stmt.QueryRow(userId).Scan(
		&log.Id,
		&log.UserId,
		&log.IpAddr,
		&log.Timestamp,
		&log.Type,
		&log.Location,
		&log.UserAgent,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get last connection: %w", err)
	}

	return &log, nil
}

func (r *LogRepositoryImpl) Create(log *models.Log) error {
	stmt, err := r.manager.GetStatement(stmtLogCreateConnection)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		log.Id,
		log.UserId,
		log.IpAddr,
		log.Timestamp,
		log.Type,
		log.Location,
		log.UserAgent,
	)

	if err != nil {
		return fmt.Errorf("failed to create connection log: %w", err)
	}

	return nil
}

func (r *LogRepositoryImpl) DeleteOld() error {
	stmt, err := r.manager.GetStatement(stmtLogDeleteOld)
	if err != nil {
		return err
	}

	// Delete logs older than 90 days
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90).UnixMilli()
	_, err = stmt.Exec(ninetyDaysAgo)
	if err != nil {
		return fmt.Errorf("failed to delete old logs: %w", err)
	}

	return nil
}
