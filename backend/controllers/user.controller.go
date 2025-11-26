package controllers

import (
	"errors"
	"net/http"
	"structured-notes/app"
	"structured-notes/models"
	"structured-notes/permissions"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUsers(c *gin.Context) (int, any)
	GetUserById(c *gin.Context) (int, any)
	GetPublicUser(c *gin.Context) (int, any)
	CreateUser(c *gin.Context) (int, any)
	UpdateUser(c *gin.Context) (int, any)
	UpdatePassword(c *gin.Context) (int, any)
	DeleteUser(c *gin.Context) (int, any)
}

func NewUserController(app *app.App) UserController {
	return &Controller{
		app:        app,
		authorizer: permissions.NewAuthorizer(app.Repos.Permission),
	}
}

func (ctr *Controller) GetUsers(c *gin.Context) (int, any) {
	users, err := ctr.app.Services.User.GetAllUsers()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, users
}

func (ctr *Controller) GetPublicUser(c *gin.Context) (int, any) {
	usernameOrEmail := c.Param("query")
	if usernameOrEmail == "" {
		return http.StatusBadRequest, errors.New("username or email is required")
	}
	users, err := ctr.app.Services.User.SearchPublicUsers(usernameOrEmail)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if len(users) == 0 {
		return http.StatusNotFound, errors.New("user not found")
	}

	return http.StatusOK, users
}

func (ctr *Controller) GetUserById(c *gin.Context) (int, any) {
	connectedUserId, connectedUserRole, err := utils.GetUserContext(c)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	targetUserId, err := utils.GetTargetId(c, c.Param("userId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if allowed, err := ctr.authorizer.CanAccessUser(connectedUserId, targetUserId, connectedUserRole); !allowed || err != nil {
		return http.StatusUnauthorized, err
	}

	result, err := ctr.app.Services.User.GetUserById(targetUserId)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, result
}

func (ctr *Controller) CreateUser(c *gin.Context) (int, any) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		return http.StatusBadRequest, err
	}

	firstname := utils.StringValue(user.Firstname)
	lastname := utils.StringValue(user.Lastname)
	avatar := utils.StringValue(user.Avatar)

	createdUser, err := ctr.app.Services.User.CreateUser(user.Username, firstname, lastname, avatar, user.Email, user.Password)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusCreated, createdUser
}

func (ctr *Controller) UpdateUser(c *gin.Context) (int, any) {
	connectedUserId, connectedUserRole, err := utils.GetUserContext(c)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	targetUserId, err := utils.GetTargetId(c, c.Param("userId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if allowed, err := ctr.authorizer.CanAccessUser(connectedUserId, targetUserId, connectedUserRole); !allowed || err != nil {
		return http.StatusUnauthorized, err
	}

	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		return http.StatusBadRequest, err
	}

	updatedUser, err := ctr.app.Services.User.UpdateUser(targetUserId, user.Firstname, user.Lastname, user.Avatar, &user.Email)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, updatedUser
}

func (ctr *Controller) UpdatePassword(c *gin.Context) (int, any) {
	connectedUserId, connectedUserRole, err := utils.GetUserContext(c)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	targetUserId, err := utils.GetTargetId(c, c.Param("userId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if allowed, err := ctr.authorizer.CanAccessUser(connectedUserId, targetUserId, connectedUserRole); !allowed || err != nil {
		return http.StatusUnauthorized, err
	}

	var payload struct {
		Password string `form:"password" json:"password"`
	}
	if err := c.ShouldBind(&payload); err != nil {
		return http.StatusBadRequest, errors.New("invalid request payload")
	}

	err = ctr.app.Services.User.UpdatePassword(targetUserId, payload.Password)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, "Password updated successfully"
}

func (ctr *Controller) DeleteUser(c *gin.Context) (int, any) {
	connectedUserId, connectedUserRole, err := utils.GetUserContext(c)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	targetUserId, err := utils.GetTargetId(c, c.Param("userId"))
	if err != nil {
		return http.StatusBadRequest, err
	}

	if allowed, err := ctr.authorizer.CanAccessUser(connectedUserId, targetUserId, connectedUserRole); !allowed || err != nil {
		return http.StatusUnauthorized, err
	}

	err = ctr.app.Services.User.DeleteUser(targetUserId, ctr.app.Services.Media)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, "User deleted successfully"
}
