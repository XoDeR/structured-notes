package permissions

import (
	"errors"
	"structured-notes/models"
	"structured-notes/repositories"
	"structured-notes/types"
)

type Authorizer interface {
	CanAccessNode(userID types.Snowflake, userRole UserRole, node *models.Node, action NodeAction) (bool, NodePermissionLevel, error)
	CanAccessUser(connectedID, targetID types.Snowflake, userRole UserRole) (bool, error)
	IsAppAdmin(userRole UserRole) bool
}

type DefaultAuthorizer struct {
	permRepo repositories.PermissionRepository
}

func NewAuthorizer(permRepo repositories.PermissionRepository) Authorizer {
	return &DefaultAuthorizer{permRepo: permRepo}
}

func (a *DefaultAuthorizer) CanAccessNode(userID types.Snowflake, userRole UserRole, node *models.Node, action NodeAction) (bool, NodePermissionLevel, error) {
	if userID == node.UserId {
		return true, PermOwner, nil
	}
	if a.IsAppAdmin(userRole) {
		return true, PermOwner, nil
	}
	hasPermission, level := a.permRepo.HasPermission(userID, node.Id, int(action.RequiredLevel()))
	if hasPermission {
		return true, NodePermissionLevel(level), nil
	}
	return false, PermNone, errors.New("unauthorized")
}

func (a *DefaultAuthorizer) CanAccessUser(connectedID, targetID types.Snowflake, userRole UserRole) (bool, error) {
	if connectedID == targetID {
		return true, nil
	}

	if a.IsAppAdmin(userRole) {
		return true, nil
	}

	return false, errors.New("unauthorized")
}

func (a *DefaultAuthorizer) IsAppAdmin(userRole UserRole) bool {
	return userRole == RoleAdministrator
}
