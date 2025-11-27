package routes

import (
	"structured-notes/app"
	"structured-notes/controllers"
	"structured-notes/middlewares"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

func Permissions(app *app.App, router *gin.RouterGroup) {
	usr := router.Group("/permissions")
	permissionsCtrl := controllers.NewPermissionsController(app)

	usr.GET("/:nodeId", middlewares.Auth(), utils.ResponseFormatter(permissionsCtrl.GetNodePermissions))
	usr.POST("", middlewares.Auth(), utils.ResponseFormatter(permissionsCtrl.CreatePermission))
	usr.PATCH("/:id", middlewares.Auth(), utils.ResponseFormatter(permissionsCtrl.UpdatePermission))
	usr.DELETE("/:id", middlewares.Auth(), utils.ResponseFormatter(permissionsCtrl.DeletePermission))
}
