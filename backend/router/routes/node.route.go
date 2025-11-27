package routes

import (
	"structured-notes/app"
	"structured-notes/controllers"
	"structured-notes/middlewares"
	"structured-notes/utils"

	"github.com/gin-gonic/gin"
)

func Nodes(app *app.App, router *gin.RouterGroup) {
	node := router.Group("/nodes")
	nodeCtrl := controllers.NewNodeController(app)

	node.GET("/public/:id", utils.ResponseFormatter(nodeCtrl.GetPublicNode))
	node.GET("/shared/:userId", middlewares.Auth(), utils.ResponseFormatter(nodeCtrl.GetSharedNodes))
	node.GET("/:userId", middlewares.Auth(), utils.ResponseFormatter(nodeCtrl.GetNodes))
	node.GET("/:userId/:id", middlewares.Auth(), utils.ResponseFormatter(nodeCtrl.GetNode))
	node.POST("", middlewares.Auth(), utils.ResponseFormatter(nodeCtrl.CreateNode))
	node.PUT("/:id", middlewares.Auth(), utils.ResponseFormatter(nodeCtrl.UpdateNode))
	node.DELETE("/:id", middlewares.Auth(), utils.ResponseFormatter(nodeCtrl.DeleteNode))
}
