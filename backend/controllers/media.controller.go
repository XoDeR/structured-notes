package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"structured-notes/app"
	"structured-notes/logger"
	"structured-notes/permissions"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

type MediaController interface {
	GetBackup(c *gin.Context) (int, any)
	UploadFile(c *gin.Context) (int, any)
	UploadAvatar(c *gin.Context) (int, any)
	DeleteUpload(c *gin.Context) (int, any)
	GetMediaFile(c *gin.Context)
}

func NewMediaController(app *app.App) MediaController {
	return &Controller{
		app:        app,
		authorizer: permissions.NewAuthorizer(app.Repos.Permission),
	}
}

func (ctr *Controller) GetBackup(c *gin.Context) (int, any) {
	userId, err := utils.GetUserIdCtx(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	link, err := ctr.app.Services.Media.CreateBackup(userId)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, gin.H{"link": link}
}

func (ctr *Controller) UploadFile(c *gin.Context) (int, any) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, errors.New("failed to get file")
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to read file")
	}

	userId, err := utils.GetUserIdCtx(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	mimeType := header.Header.Get("Content-Type")
	node, err := ctr.app.Services.Media.UploadFile(
		header.Filename,
		header.Size,
		fileContent,
		mimeType,
		userId,
		ctr.app.Config.Media.MaxSize,
		ctr.app.Config.Media.MaxUploadsSize,
		ctr.app.Config.Media.SupportedTypes,
	)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, node
}

func (ctr *Controller) UploadAvatar(c *gin.Context) (int, any) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return http.StatusBadRequest, errors.New("failed to get file")
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to read file")
	}

	userId, err := utils.GetUserIdCtx(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	mimeType := header.Header.Get("Content-Type")
	err = ctr.app.Services.Media.UploadAvatar(
		header.Filename,
		header.Size,
		fileContent,
		mimeType,
		userId,
		ctr.app.Config.Media.MaxSize,
		ctr.app.Config.Media.SupportedTypesImages,
	)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, "Avatar uploaded successfully"
}

func (ctr *Controller) DeleteUpload(c *gin.Context) (int, any) {
	connectedUserId, connectedUserRole, err := utils.GetUserContext(c)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	nodeTargetId, err := utils.GetTargetId(c, c.Param("id"))
	if err != nil {
		return http.StatusBadRequest, errors.New("invalid id format")
	}

	err = ctr.app.Services.Media.DeleteUpload(nodeTargetId, connectedUserId, connectedUserRole, ctr.authorizer)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	return http.StatusOK, "Media deleted successfully"
}

func (ctr *Controller) GetMediaFile(c *gin.Context) {
	logger.Info("mediaUploads path hit")

	connectedUserId, connectedUserRole, err := utils.GetUserContext(c)
	if err != nil {
		logger.Info("GetMediaFile: connectedUserId: " + strconv.FormatUint(uint64(connectedUserId), 10))
		logger.Info("GetMediaFile: connectedUserRole: " + strconv.FormatInt(int64(connectedUserRole), 10))
		return
	}

	userTargetId, err := utils.GetTargetId(c, c.Param("userId"))
	if err != nil {
		logger.Info("GetMediaFile: userTargetId: " + strconv.FormatUint(uint64(userTargetId), 10))
		return
	}

	nodeTargetId, ext, err := utils.GetMediaFilenameParts(c, c.Param("nameAndExt"))
	if err != nil {
		logger.Info("GetMediaFile: nodeTargetId: " + strconv.FormatUint(uint64(nodeTargetId), 10))
		logger.Info("GetMediaFile: ext: " + ext)
		return
	}

	// checks if the user has permissions for the media
	err = ctr.app.Services.Media.GetMediaFile(nodeTargetId, userTargetId, connectedUserId, connectedUserRole, ctr.authorizer)
	if err != nil {
		logger.Info("GetMediaFile")
		return
	}

	userIdStr := fmt.Sprintf("%d", userTargetId)
	nodeIdStr := fmt.Sprintf("%d", nodeTargetId)
	nodeIdAndExtStr := nodeIdStr + ext

	fullMediaFilePath := filepath.Join("media", userIdStr, nodeIdAndExtStr)

	c.Header("Cross-Origin-Resource-Policy", "cross-origin")

	c.File(fullMediaFilePath)

	logger.Info("File attached")
}
