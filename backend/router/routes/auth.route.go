package routes

import (
	"structured-notes/app"
	"structured-notes/controllers"
	"structured-notes/middlewares"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

func Auth(app *app.App, router *gin.RouterGroup) {
	auth := router.Group("/auth")

	authCtrl := controllers.NewAuthController(app)
	auth.POST("", utils.ResponseFormatter(authCtrl.Login))
	auth.POST("/refresh", utils.ResponseFormatter(authCtrl.RefreshSession))
	auth.POST("/request-reset", utils.ResponseFormatter(authCtrl.RequestResetPassword))
	auth.POST("/reset-password", utils.ResponseFormatter(authCtrl.ResetPassword))
	auth.POST("/logout", utils.ResponseFormatter(authCtrl.Logout))
	auth.POST("/logout/all", middlewares.Auth(), utils.ResponseFormatter(authCtrl.LogoutAllDevices))
}
