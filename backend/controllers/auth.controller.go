package controllers

import (
	"errors"
	"net/http"
	"os"
	"structured-notes/app"
	"structured-notes/logger"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context) (int, any)
	RefreshSession(c *gin.Context) (int, any)
	RequestResetPassword(c *gin.Context) (int, any)
	ResetPassword(c *gin.Context) (int, any)
	Logout(c *gin.Context) (int, any)
	LogoutAllDevices(c *gin.Context) (int, any)
}

func NewAuthController(app *app.App) AuthController {
	deleteOldSessionsAndLogs(app)
	return &Controller{app: app}
}

type AuthClaims struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func (ctr *Controller) Login(c *gin.Context) (int, any) {
	var authClaims AuthClaims
	if err := c.ShouldBind(&authClaims); err != nil {
		return http.StatusBadRequest, err
	}

	user, session, err := ctr.app.Services.Auth.Login(authClaims.Username, authClaims.Password, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		return http.StatusUnauthorized, err
	}

	tokenString, err := ctr.app.Services.Auth.SignAccessToken(user, ctr.app.Config.Auth.AccessTokenExpiry)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to sign token")
	}

	secure := shouldUseSecureCookies()
	c.SetCookie("Authorization", tokenString, ctr.app.Config.Auth.AccessTokenExpiry, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)
	c.SetCookie("RefreshToken", session.RefreshToken, ctr.app.Config.Auth.RefreshTokenExpiry, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)

	return http.StatusOK, user
}

func (ctr *Controller) RefreshSession(c *gin.Context) (int, any) {
	refreshToken, err := c.Cookie("RefreshToken")
	if err != nil {
		return http.StatusUnauthorized, errors.New("no refresh token provided")
	}

	user, session, err := ctr.app.Services.Auth.RefreshSession(refreshToken)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	tokenString, err := ctr.app.Services.Auth.SignAccessToken(user, ctr.app.Config.Auth.AccessTokenExpiry)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to sign token")
	}

	secure := shouldUseSecureCookies()
	c.SetCookie("Authorization", tokenString, ctr.app.Config.Auth.AccessTokenExpiry, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)
	c.SetCookie("RefreshToken", session.RefreshToken, ctr.app.Config.Auth.RefreshTokenExpiry, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)

	return http.StatusOK, "Session refreshed successfully."
}

func (ctr *Controller) Logout(c *gin.Context) (int, any) {
	refreshToken, err := c.Cookie("RefreshToken")
	if err != nil {
		return http.StatusUnauthorized, errors.New("no refresh token provided")
	}

	if err := ctr.app.Services.Auth.Logout(refreshToken); err != nil {
		return http.StatusUnauthorized, err
	}

	secure := shouldUseSecureCookies()
	c.SetCookie("Authorization", "", -1, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)
	c.SetCookie("RefreshToken", "", -1, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)
	return http.StatusOK, "Logged out successfully."
}

func (ctr *Controller) LogoutAllDevices(c *gin.Context) (int, any) {
	userId, err := utils.GetUserIdCtx(c)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if err := ctr.app.Services.Auth.LogoutAllDevices(userId); err != nil {
		return http.StatusInternalServerError, errors.New("failed to delete sessions")
	}

	secure := shouldUseSecureCookies()
	c.SetCookie("Authorization", "", -1, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)
	c.SetCookie("RefreshToken", "", -1, "/", os.Getenv("COOKIE_DOMAIN"), secure, true)
	return http.StatusOK, "Logged out from all devices successfully."
}

func (ctr *Controller) RequestResetPassword(c *gin.Context) (int, any) {
	var data struct {
		User string `json:"username" binding:"required"`
	}
	if err := c.ShouldBind(&data); err != nil {
		return http.StatusBadRequest, err
	}

	ctr.app.Services.Auth.RequestPasswordReset(data.User)
	return http.StatusOK, "Job done."
}

func (ctr *Controller) ResetPassword(c *gin.Context) (int, any) {
	var data struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&data); err != nil {
		return http.StatusBadRequest, err
	}

	if err := ctr.app.Services.Auth.ResetPassword(data.Token, data.Password); err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, "Password reset successfully."
}

func deleteOldSessionsAndLogs(app *app.App) {
	err := app.Services.Log.DeleteOldLogs()
	if err != nil {
		logger.Error("Error deleting old logs: " + err.Error())
	} else {
		logger.Success("Old logs deleted successfully.")
	}
	err = app.Services.Session.DeleteOldSessions()
	if err != nil {
		logger.Error("Error deleting old sessions: " + err.Error())
	} else {
		logger.Success("Old sessions deleted successfully.")
	}
}

func shouldUseSecureCookies() bool {
	value := os.Getenv("ALLOW_UNSECURE")
	return !(value == "true" || value == "1")
}
