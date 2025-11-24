package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"structured-notes/models"
	"structured-notes/types"
)

type NodeRepository interface {
	GetAll(userId types.Snowflake) ([]*models.Node, error)
	GetShared(userId types.Snowflake) ([]*models.Node, error)
	GetAllForBackup(userId types.Snowflake) ([]*models.Node, error)
	GetByID(nodeId types.Snowflake) (*models.Node, error)
	GetPublic(nodeId types.Snowflake) (*models.Node, error)
	GetUserUploadsSize(userId types.Snowflake) (int64, error)
	Create(node *models.Node) error
	Update(node *models.Node) error
	Delete(nodeId types.Snowflake) error
}

type NodeRepositoryImpl struct {
	db      *sql.DB
	manager *RepositoryManager
}

const (
	stmNodeGetAll              = "node_get_all"
	stmNodeGetShared           = "node_get_shared"
	stmNodeGetAllForBackup     = "node_get_all_backup"
	stmtNodeGetByID            = "node_get_by_id"
	stmtNodeGetPublic          = "node_get_public"
	stmtNodeGetUserUploadsSize = "node_get_user_uploads_size"
	stmtNodeCreate             = "node_create"
	stmtNodeUpdate             = "node_update"
	stmtNodeDelete             = "node_delete"
)

func NewNodeRepository(db *sql.DB, manager *RepositoryManager) (NodeRepository, error) {
	repo := &NodeRepositoryImpl{
		db:      db,
		manager: manager,
	}

	if err := repo.prepareStatements(); err != nil {
		return nil, fmt.Errorf("failed to prepare node statements: %w", err)
	}

	return repo, nil
}

func (r *NodeRepositoryImpl) prepareStatements() error {
	statements := map[string]string{

		stmNodeGetAll: `
		WITH RECURSIVE user_nodes AS (
		SELECT n.id, n.user_id, n.parent_id, n.name, n.description, n.tags, n.role, n.color, n.icon, n.theme,
				   n.accessibility, n.access, n.display, n.order, n.size, n.metadata, n.created_timestamp, n.updated_timestamp
		FROM nodes n
		WHERE n.user_id = ?

		UNION

		SELECT c.id, c.user_id, c.parent_id, c.name, c.description, c.tags, c.role, c.color, c.icon, c.theme,
				   c.accessibility, c.access, c.display, c.order, c.size, c.metadata, c.created_timestamp, c.updated_timestamp
		FROM nodes c
		JOIN user_nodes un ON un.id = c.parent_id)
		SELECT * FROM user_nodes ORDER BY role, 'order' DESC, name;`,

		stmNodeGetShared: `
		WITH RECURSIVE shared_nodes AS (
		    SELECT n.id, n.user_id, n.parent_id, n.name, n.description, n.tags, n.role, n.color, n.icon, n.theme,
		           n.accessibility, n.access, n.display, n.order, n.size, n.metadata, n.created_timestamp, n.updated_timestamp
		    FROM nodes n
		    JOIN permissions p ON p.node_id = n.id
		    WHERE p.user_id = ?

		    UNION

		    SELECT c.id, c.user_id, c.parent_id, c.name, c.description, c.tags, c.role, c.color, c.icon, c.theme,
		           c.accessibility, c.access, c.display, c.order, c.size, c.metadata, c.created_timestamp, c.updated_timestamp
		    FROM nodes c
		    JOIN shared_nodes an ON an.id = c.parent_id
		)
		SELECT * FROM shared_nodes;`,

		stmNodeGetAllForBackup: `
		SELECT id, user_id, parent_id, name, description, tags, role, color, icon, thumbnail, theme, 
		       accessibility, access, display, ` + "`order`" + `, content, content_compiled, size, metadata, 
		       created_timestamp, updated_timestamp 
		FROM nodes 
		WHERE user_id = ?`,

		stmtNodeGetByID: `
			SELECT id, user_id, parent_id, name, description, tags, role, color, icon, thumbnail, theme, 
			       accessibility, access, display, ` + "`order`" + `, content, content_compiled, size, metadata, 
			       created_timestamp, updated_timestamp 
			FROM nodes 
			WHERE id = ?`,

		stmtNodeGetPublic: `
			SELECT id, user_id, parent_id, name, description, tags, role, color, icon, thumbnail, theme, 
			       accessibility, access, display, ` + "`order`" + `, content, content_compiled, size, metadata, 
			       created_timestamp, updated_timestamp 
			FROM nodes 
			WHERE id = ? AND accessibility = 3`,

		stmtNodeGetUserUploadsSize: `
			SELECT COALESCE(SUM(size), 0) 
			FROM nodes 
			WHERE user_id = ?`,

		stmtNodeCreate: `
			INSERT INTO nodes (id, user_id, parent_id, name, description, tags, role, color, icon, thumbnail, theme, 
			                   accessibility, access, display, ` + "`order`" + `, content, content_compiled, size, metadata, 
			                   created_timestamp, updated_timestamp) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,

		stmtNodeUpdate: `
			UPDATE nodes 
			SET parent_id = ?, user_id = ?, name = ?, description = ?, tags = ?, role = ?, color = ?, 
			    icon = ?, thumbnail = ?, theme = ?, accessibility = ?, access = ?, display = ?, ` + "`order`" + ` = ?, 
			    content = ?, content_compiled = ?, metadata = ?, updated_timestamp = ? 
			WHERE id = ?`,

		stmtNodeDelete: `
			DELETE FROM nodes 
			WHERE id = ?`,
	}

	for key, query := range statements {
		if _, err := r.manager.PrepareStatement(key, query); err != nil {
			return err
		}
	}

	return nil
}

func (r *NodeRepositoryImpl) scanNode(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Node, error) {
	var node models.Node
	err := scanner.Scan(
		&node.Id,
		&node.UserId,
		&node.ParentId,
		&node.Name,
		&node.Description,
		&node.Tags,
		&node.Role,
		&node.Color,
		&node.Icon,
		&node.Thumbnail,
		&node.Theme,
		&node.Accessibility,
		&node.Access,
		&node.Display,
		&node.Order,
		&node.Content,
		&node.ContentCompiled,
		&node.Size,
		&node.Metadata,
		&node.CreatedTimestamp,
		&node.UpdatedTimestamp,
	)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NodeRepositoryImpl) scanNodePartial(scanner interface {
	Scan(dest ...interface{}) error
}) (*models.Node, error) {
	var node models.Node
	err := scanner.Scan(
		&node.Id,
		&node.UserId,
		&node.ParentId,
		&node.Name,
		&node.Description,
		&node.Tags,
		&node.Role,
		&node.Color,
		&node.Icon,
		&node.Theme,
		&node.Accessibility,
		&node.Access,
		&node.Display,
		&node.Order,
		&node.Size,
		&node.Metadata,
		&node.CreatedTimestamp,
		&node.UpdatedTimestamp,
	)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NodeRepositoryImpl) GetAll(userId types.Snowflake) ([]*models.Node, error) {
	stmt, err := r.manager.GetStatement("node_get_all")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query user nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]*models.Node, 0)
	for rows.Next() {
		node, err := r.scanNodePartial(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nodes: %w", err)
	}

	return nodes, nil
}

func (r *NodeRepositoryImpl) GetShared(userId types.Snowflake) ([]*models.Node, error) {
	stmt, err := r.manager.GetStatement("node_get_shared")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query shared nodes: %w", err)
	}
	defer rows.Close()

	nodeMap := make(map[types.Snowflake]*models.Node)
	nodes := make([]*models.Node, 0)

	for rows.Next() {
		node, err := r.scanNodePartial(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		node.Permissions = []*models.Permission{}
		nodes = append(nodes, node)
		nodeMap[node.Id] = node
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nodes: %w", err)
	}

	if len(nodes) > 0 {
		if err := r.loadPermissionsForNodes(userId, nodeMap); err != nil {
			return nil, fmt.Errorf("failed to load permissions: %w", err)
		}
	}

	return nodes, nil
}

func (r *NodeRepositoryImpl) loadPermissionsForNodes(userId types.Snowflake, nodeMap map[types.Snowflake]*models.Node) error {
	if len(nodeMap) == 0 {
		return nil
	}

	nodeIDs := make([]interface{}, 0, len(nodeMap))
	placeholders := make([]string, 0, len(nodeMap))
	for nodeId := range nodeMap {
		nodeIDs = append(nodeIDs, nodeId)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(`
		SELECT id, node_id, user_id, permission, created_timestamp 
		FROM permissions 
		WHERE user_id = ? AND node_id IN (%s)`,
		strings.Join(placeholders, ","))

	// Prepare this query (it's dynamic so we prepare it on-the-fly)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare permissions query: %w", err)
	}
	defer stmt.Close()

	args := append([]interface{}{userId}, nodeIDs...)
	rows, err := stmt.Query(args...)
	if err != nil {
		return fmt.Errorf("failed to query permissions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Permission
		if err := rows.Scan(&p.Id, &p.NodeId, &p.UserId, &p.Permission, &p.CreatedTimestamp); err != nil {
			return fmt.Errorf("failed to scan permission: %w", err)
		}

		if node, ok := nodeMap[p.NodeId]; ok {
			node.Permissions = append(node.Permissions, &p)
		}
	}

	return rows.Err()
}

func (r *NodeRepositoryImpl) GetAllForBackup(userId types.Snowflake) ([]*models.Node, error) {
	stmt, err := r.manager.GetStatement("node_get_all_backup")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes for backup: %w", err)
	}
	defer rows.Close()

	nodes := make([]*models.Node, 0)
	for rows.Next() {
		node, err := r.scanNode(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating nodes: %w", err)
	}

	return nodes, nil
}

func (r *NodeRepositoryImpl) GetByID(nodeId types.Snowflake) (*models.Node, error) {
	stmt, err := r.manager.GetStatement(stmtNodeGetByID)
	if err != nil {
		return nil, err
	}

	node, err := r.scanNode(stmt.QueryRow(nodeId))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get node by id: %w", err)
	}

	return node, nil
}

func (r *NodeRepositoryImpl) GetPublic(nodeId types.Snowflake) (*models.Node, error) {
	stmt, err := r.manager.GetStatement(stmtNodeGetPublic)
	if err != nil {
		return nil, err
	}

	node, err := r.scanNode(stmt.QueryRow(nodeId))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get public node: %w", err)
	}

	return node, nil
}

func (r *NodeRepositoryImpl) GetUserUploadsSize(userId types.Snowflake) (int64, error) {
	stmt, err := r.manager.GetStatement(stmtNodeGetUserUploadsSize)
	if err != nil {
		return 0, err
	}

	var totalSize int64
	err = stmt.QueryRow(userId).Scan(&totalSize)
	if err != nil {
		return 0, fmt.Errorf("failed to get user uploads size: %w", err)
	}

	return totalSize, nil
}

func (r *NodeRepositoryImpl) Create(node *models.Node) error {
	stmt, err := r.manager.GetStatement(stmtNodeCreate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		node.Id,
		node.UserId,
		node.ParentId,
		node.Name,
		node.Description,
		node.Tags,
		node.Role,
		node.Color,
		node.Icon,
		node.Thumbnail,
		node.Theme,
		node.Accessibility,
		node.Access,
		node.Display,
		node.Order,
		node.Content,
		node.ContentCompiled,
		node.Size,
		node.Metadata,
		node.CreatedTimestamp,
		node.UpdatedTimestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}

	return nil
}

func (r *NodeRepositoryImpl) Update(node *models.Node) error {
	stmt, err := r.manager.GetStatement(stmtNodeUpdate)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		node.ParentId,
		node.UserId,
		node.Name,
		node.Description,
		node.Tags,
		node.Role,
		node.Color,
		node.Icon,
		node.Thumbnail,
		node.Theme,
		node.Accessibility,
		node.Access,
		node.Display,
		node.Order,
		node.Content,
		node.ContentCompiled,
		node.Metadata,
		node.UpdatedTimestamp,
		node.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	return nil
}

func (r *NodeRepositoryImpl) Delete(nodeId types.Snowflake) error {
	stmt, err := r.manager.GetStatement(stmtNodeDelete)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(nodeId)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	return nil
}
