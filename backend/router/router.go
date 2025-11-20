package router

import (
	"structured-notes/app"

	"github.com/gin-gonic/gin"
)

func InitRouter(app *app.App) *gin.Engine {
	router := gin.New()

	return router
}
