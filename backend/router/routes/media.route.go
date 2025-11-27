package routes

import (
	"structured-notes/app"
	"structured-notes/controllers"
	"structured-notes/middlewares"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

func Uploads(app *app.App, router *gin.RouterGroup) {
	media := router.Group("/media")
	mediaCtrl := controllers.NewMediaController(app)

	media.Use(middlewares.Auth())
	media.GET(("/backup"), utils.ResponseFormatter(mediaCtrl.GetBackup))
	media.POST("", utils.ResponseFormatter(mediaCtrl.UploadFile))
	media.POST("/avatar", utils.ResponseFormatter(mediaCtrl.UploadAvatar))
	media.DELETE("/:id", utils.ResponseFormatter(mediaCtrl.DeleteUpload))
}
