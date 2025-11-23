package repositories

import (
	"database/sql"
	"fmt"
	"structured-notes/models"
	"structured-notes/types"
)

type PermissionRepository interface {
	GetByID(permissionId types.Snowflake) (*models.Permission, error)
	GetByNode(nodeId types.Snowflake) ([]*models.Permission, error)
	GetByNodeAndUser(nodeId types.Snowflake, userId types.Snowflake) (*models.Permission, error)
	HasPermission(userId, nodeId types.Snowflake, required int) (bool, int)
	Create(permission *models.Permission) (*models.Permission, error)
	Update(permission *models.Permission) error
	Delete(permissionId types.Snowflake) error
}

type PermissionRepositoryImpl struct {
	db      *sql.DB
	manager *RepositoryManager
}

const (
	stmtPermissionGetByID          = "permission_get_by_id"
	stmtPermissionGetByNode        = "permission_get_by_node"
	stmtPermissionGetByNodeAndUser = "permission_get_by_node_and_user"
	stmtPermissionCreate           = "permission_create"
	stmtPermissionUpdate           = "permission_update"
	stmtPermissionDelete           = "permission_delete"
)

func NewPermissionRepository(db *sql.DB, manager *RepositoryManager) (PermissionRepository, error) {
	repo := &PermissionRepositoryImpl{
		db:      db,
		manager: manager,
	}

	if err := repo.prepareStatements(); err != nil {
		return nil, fmt.Errorf("failed to prepare permission statements: %w", err)
	}

	return repo, nil
}

func (r *PermissionRepositoryImpl) prepareStatements() error {
	statements := map[string]string{
		stmtPermissionGetByID: `
			SELECT id, node_id, user_id, permission, created_timestamp
			FROM permissions
			WHERE id = ?`,

		stmtPermissionGetByNode: `
			SELECT id, node_id, user_id, permission, created_timestamp
			FROM permissions
			WHERE node_id = ?`,

		stmtPermissionGetByNodeAndUser: `
			SELECT id, node_id, user_id, permission, created_timestamp
			FROM permissions
			WHERE node_id = ? AND user_id = ?`,

		stmtPermissionCreate: `
			INSERT INTO permissions (id, node_id, user_id, permission, created_timestamp)
			VALUES (?, ?, ?, ?, ?)`,

		stmtPermissionUpdate: `
			UPDATE permissions
			SET permission = ?
			WHERE id = ?`,

		stmtPermissionDelete: `
			DELETE FROM permissions
			WHERE id = ?`,
	}

	for key, query := range statements {
		if _, err := r.manager.PrepareStatement(key, query); err != nil {
			return err
		}
	}

	return nil
}

func (r *PermissionRepositoryImpl) GetByID(permissionId types.Snowflake) (*models.Permission, error) {
	stmt, err := r.manager.GetStatement(stmtPermissionGetByID)
	if err != nil {
		return nil, err
	}

	var perm models.Permission
	err = stmt.QueryRow(permissionId).Scan(
		&perm.Id,
		&perm.NodeId,
		&perm.UserId,
		&perm.Permission,
		&perm.CreatedTimestamp,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &perm, nil
}

func (r *PermissionRepositoryImpl) GetByNode(nodeId types.Snowflake) ([]*models.Permission, error) {
	stmt, err := r.manager.GetStatement(stmtPermissionGetByNode)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(nodeId)
	if err != nil {
		return nil, fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	permissions := make([]*models.Permission, 0)
	for rows.Next() {
		var perm models.Permission
		err := rows.Scan(
			&perm.Id,
			&perm.NodeId,
			&perm.UserId,
			&perm.Permission,
			&perm.CreatedTimestamp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, &perm)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating permissions: %w", err)
	}

	return permissions, nil
}

func (r *PermissionRepositoryImpl) GetByNodeAndUser(nodeId types.Snowflake, userId types.Snowflake) (*models.Permission, error) {
	stmt, err := r.manager.GetStatement(stmtPermissionGetByNodeAndUser)
	if err != nil {
		return nil, err
	}

	var perm models.Permission
	err = stmt.QueryRow(nodeId, userId).Scan(
		&perm.Id,
		&perm.NodeId,
		&perm.UserId,
		&perm.Permission,
		&perm.CreatedTimestamp,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return &perm, nil
}

func (r *PermissionRepositoryImpl) Create(permission *models.Permission) (*models.Permission, error) {
	stmt, err := r.manager.GetStatement(stmtPermissionCreate)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(
		permission.Id,
		permission.NodeId,
		permission.UserId,
		permission.Permission,
		permission.CreatedTimestamp,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}

func (r *PermissionRepositoryImpl) Update(permission *models.Permission) error {
	stmt, err := r.manager.GetStatement(stmtPermissionUpdate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		permission.Permission,
		permission.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}

func (r *PermissionRepositoryImpl) Delete(permissionId types.Snowflake) error {
	stmt, err := r.manager.GetStatement(stmtPermissionDelete)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(permissionId)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

func (r *PermissionRepositoryImpl) HasPermission(userId, nodeId types.Snowflake, required int) (bool, int) {
	var perm sql.NullInt32
	r.db.QueryRow(`
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_id
			FROM nodes
			WHERE id = ?

			UNION ALL

			SELECT n.id, n.parent_id
			FROM nodes n
			INNER JOIN ancestors a ON a.parent_id = n.id
		)
		SELECT MAX(p.permission)
		FROM permissions p
		INNER JOIN ancestors an ON an.id = p.node_id
		WHERE p.user_id = ?
	`, nodeId, userId).Scan(&perm)

	// Explicit permission is sufficient
	if perm.Valid && perm.Int32 >= int32(required) {
		return true, int(perm.Int32)
	}

	// Not enough permissions ==> check if owner of an ancestor
	var owns int
	err := r.db.QueryRow(`
		WITH RECURSIVE ancestors AS (
			SELECT id, parent_id, user_id
			FROM nodes
			WHERE id = ?

			UNION ALL

			SELECT n.id, n.parent_id, n.user_id
			FROM nodes n
			INNER JOIN ancestors a ON a.parent_id = n.id
		)
		SELECT 1
		FROM ancestors
		WHERE user_id = ?
		LIMIT 1
	`, nodeId, userId).Scan(&owns)
	if err == nil && owns == 1 {
		return true, 4
	}

	return false, 0
}
