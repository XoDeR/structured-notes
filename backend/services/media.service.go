package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"structured-notes/logger"
	"structured-notes/models"
	"structured-notes/permissions"
	"structured-notes/repositories"
	"structured-notes/types"
	"structured-notes/utils"
)

type MediaService interface {
	CreateBackup(userId types.Snowflake) (string, error)
	UploadFile(filename string, fileSize int64, fileContent []byte, mimeType string, userId types.Snowflake, maxSize, maxUploadsSize float64, supportedTypes []string) (*models.Node, error)
	UploadAvatar(filename string, fileSize int64, fileContent []byte, mimeType string, userId types.Snowflake, maxSize float64, supportedTypes []string) error
	DeleteUpload(nodeId types.Snowflake, connectedUserId types.Snowflake, connectedUserRole permissions.UserRole, authorizer permissions.Authorizer) error
	DeleteAllFromUser(userId types.Snowflake) error
	GetMediaFile(nodeId types.Snowflake, userId types.Snowflake, connectedUserId types.Snowflake, connectedUserRole permissions.UserRole, authorizer permissions.Authorizer) error
}

type mediaService struct {
	nodeRepo  repositories.NodeRepository
	snowflake *utils.Snowflake
}

func NewMediaService(nodeRepo repositories.NodeRepository, snowflake *utils.Snowflake) MediaService {
	return &mediaService{
		nodeRepo:  nodeRepo,
		snowflake: snowflake,
	}
}

func (s *mediaService) CreateBackup(userId types.Snowflake) (string, error) {
	nodes, err := s.nodeRepo.GetAllForBackup(userId)
	if err != nil {
		return "", err
	}

	backup := map[string]interface{}{
		"nodes": nodes,
	}
	jsonString, err := json.Marshal(backup)
	if err != nil {
		return "", err
	}

	objectName := fmt.Sprintf("%d/backups/backup.json", userId)

	// NOT IMPLEMENTED
	logger.Info("NOT IMPLEMENTED")
	n := 100
	if len(jsonString) < n {
		n = len(jsonString)
	}
	strToLog := string(jsonString[:n])
	logger.Info("backup" + strToLog)

	return objectName, nil
}

func (s *mediaService) UploadFile(filename string, fileSize int64, fileContent []byte, mimeType string, userId types.Snowflake, maxSize, maxUploadsSize float64, supportedTypes []string) (*models.Node, error) {
	if fileSize > int64(maxSize) {
		return nil, errors.New("file size exceeds the limit")
	}

	if !slices.Contains(supportedTypes, mimeType) {
		return nil, errors.New("file type not supported")
	}

	totalSize, err := s.nodeRepo.GetUserUploadsSize(userId)
	if err != nil {
		return nil, err
	}
	if totalSize+fileSize > int64(maxUploadsSize) {
		return nil, errors.New("total size of uploads exceeds the limit")
	}

	id := s.snowflake.Generate()
	ext := filepath.Ext(filename)
	transformedPath := fmt.Sprintf("%d%s", id, ext)
	metadata := types.JSONB{
		"filetype":         mimeType,
		"original_path":    filename,
		"transformed_path": transformedPath,
	}

	name := filename
	if len(name) > 50 {
		name = name[:50]
	}
	accessibility := utils.IntPtr(1)

	node := &models.Node{
		Id:              id,
		UserId:          userId,
		ParentId:        nil,
		Name:            name,
		Role:            4,
		Accessibility:   accessibility,
		Access:          0,
		Size:            &fileSize,
		Content:         &filename,
		ContentCompiled: &transformedPath,
		Metadata:        &metadata,
	}

	if err := s.nodeRepo.Create(node); err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("%d/%d%s", userId, id, ext)

	// NOT IMPLEMENTED
	logger.Info("NOT IMPLEMENTED")
	logger.Info("filename" + objectName)

	return node, nil
}

func (s *mediaService) UploadAvatar(filename string, fileSize int64, fileContent []byte, mimeType string, userId types.Snowflake, maxSize float64, supportedTypes []string) error {
	if fileSize > int64(maxSize) {
		return errors.New("file size exceeds the limit")
	}

	if !slices.Contains(supportedTypes, mimeType) {
		return errors.New("file type not supported")
	}

	objectName := fmt.Sprintf("%d/avatar", userId)
	// NOT IMPLEMENTED
	logger.Info("NOT IMPLEMENTED")
	logger.Info("filename" + objectName)

	return nil
}

func (s *mediaService) DeleteUpload(nodeId types.Snowflake, connectedUserId types.Snowflake, connectedUserRole permissions.UserRole, authorizer permissions.Authorizer) error {
	node, err := s.nodeRepo.GetByID(nodeId)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	if allowed, err := authorizer.CanAccessUser(connectedUserId, node.UserId, connectedUserRole); !allowed || err != nil {
		return errors.New("unauthorized")
	}

	prefix := fmt.Sprintf("%d/%d", node.UserId, node.Id)
	// NOT IMPLEMENTED
	logger.Info("NOT IMPLEMENTED")
	logger.Info("filename prefix" + prefix)

	return s.nodeRepo.Delete(nodeId)
}

func (s *mediaService) DeleteAllFromUser(userId types.Snowflake) error {
	prefix := fmt.Sprintf("%d/", userId)
	// NOT IMPLEMENTED
	logger.Info("NOT IMPLEMENTED")
	logger.Info("filename prefix" + prefix)

	return nil
}

func (s *mediaService) GetMediaFile(nodeId types.Snowflake, userId types.Snowflake, connectedUserId types.Snowflake, connectedUserRole permissions.UserRole, authorizer permissions.Authorizer) error {
	node, err := s.nodeRepo.GetByID(nodeId)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	if allowed, err := authorizer.CanAccessUser(connectedUserId, node.UserId, connectedUserRole); !allowed || err != nil {
		return errors.New("unauthorized")
	}

	return nil
}
