package services

import (
	"fmt"
	"structured-notes/logger"
	"structured-notes/repositories"
	"structured-notes/utils"
)

type ServiceManager struct {
	Auth        AuthService
	User        UserService
	Node        NodeService
	Permission  PermissionService
	Log         LogService
	Session     SessionService
	Media       MediaService
	initialized bool
}

func NewServiceManager(repos *repositories.RepositoryManager, snowflake *utils.Snowflake) (*ServiceManager, error) {
	sm := &ServiceManager{}

	if err := sm.initializeServices(repos, snowflake); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	sm.initialized = true
	logger.Success("Service manager initialized successfully")
	return sm, nil
}

func (sm *ServiceManager) initializeServices(repos *repositories.RepositoryManager, snowflake *utils.Snowflake) error {
	sm.Auth = NewAuthService(repos.User, repos.Session, repos.Log, snowflake)
	sm.User = NewUserService(repos.User, repos.Log, snowflake)
	sm.Node = NewNodeService(repos.Node, repos.Permission, snowflake)
	sm.Permission = NewPermissionService(repos.Permission, repos.Node, snowflake)
	sm.Log = NewLogService(repos.Log, snowflake)
	sm.Session = NewSessionService(repos.Session)
	sm.Media = NewMediaService(repos.Node, snowflake)

	return nil
}

func (sm *ServiceManager) IsInitialized() bool {
	return sm.initialized
}
