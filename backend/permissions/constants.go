package permissions

// UserRole
type UserRole int

const (
	RoleNone          UserRole = 0
	RoleAdministrator UserRole = 1 << 1
	RoleManager       UserRole = 1 << 2
	RoleModerator     UserRole = 1 << 3
)

// NodePermissionLevel
type NodePermissionLevel int

const (
	PermNone NodePermissionLevel = iota
	PermRead
	PermWrite
	PermAdmin
	PermOwner // Full, including managing permissions
)

// NodeAction
type NodeAction int

const (
	ActionRead NodeAction = iota + 1
	ActionUpdate
	ActionDelete
	ActionShare
	ActionManagePermissions
)

func (nodeAction NodeAction) RequiredLevel() NodePermissionLevel {
	switch nodeAction {
	case ActionRead:
		return PermRead
	case ActionUpdate:
		return PermWrite
	case ActionDelete, ActionShare:
		return PermAdmin
	case ActionManagePermissions:
		return PermOwner
	default:
		return PermOwner
	}
}
