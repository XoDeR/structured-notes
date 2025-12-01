package routes

import (
	"structured-notes/app"
	"structured-notes/controllers"
	"structured-notes/middlewares"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

func Uploads(app *app.App, mainGroup *gin.RouterGroup, mediaGroup *gin.RouterGroup) {
	// /api/media
	media := mainGroup.Group("/media")
	mediaCtrl := controllers.NewMediaController(app)

	media.Use(middlewares.Auth())
	media.GET(("/backup"), utils.ResponseFormatter(mediaCtrl.GetBackup))
	media.POST("", utils.ResponseFormatter(mediaCtrl.UploadFile))
	media.POST("/avatar", utils.ResponseFormatter(mediaCtrl.UploadAvatar))
	media.DELETE("/:id", utils.ResponseFormatter(mediaCtrl.DeleteUpload))

	// /media
	// Processes GET from for example <img src="[serverUrl]/media/[userId]/[nodeId].png">
	mediaUploads := mediaGroup
	mediaUploads.Use(middlewares.Auth())
	mediaUploads.GET("/:userId/:nameAndExt", mediaCtrl.GetMediaFile)
}
