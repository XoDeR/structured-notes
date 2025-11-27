package routes

import (
	"structured-notes/app"
	"structured-notes/controllers"
	"structured-notes/middlewares"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

func Users(app *app.App, router *gin.RouterGroup) {
	usr := router.Group("/users")
	usrCtrl := controllers.NewUserController(app)

	usr.GET("", middlewares.Auth(), middlewares.Admin(), utils.ResponseFormatter(usrCtrl.GetUsers))
	usr.GET("/:userId", middlewares.Auth(), utils.ResponseFormatter(usrCtrl.GetUserById))
	usr.GET("/public/:query", utils.ResponseFormatter(usrCtrl.GetPublicUser))
	usr.POST("", utils.ResponseFormatter(usrCtrl.CreateUser))
	usr.PATCH("/:userId", middlewares.Auth(), utils.ResponseFormatter(usrCtrl.UpdateUser))
	usr.PATCH("/:userId/password", middlewares.Auth(), utils.ResponseFormatter(usrCtrl.UpdatePassword))
	usr.DELETE("/:userId", middlewares.Auth(), utils.ResponseFormatter(usrCtrl.DeleteUser))
}
