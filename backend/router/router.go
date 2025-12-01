package router

import (
	"os"
	"structured-notes/app"
	"structured-notes/router/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(app *app.App) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("DOMAIN_CLIENT")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept"},
		AllowCredentials: true,
	}))

	mainGroup := router.Group("/api")
	mediaGroup := router.Group("/media")
	routes.Users(app, mainGroup)
	routes.Auth(app, mainGroup)
	routes.Uploads(app, mainGroup, mediaGroup)
	routes.Nodes(app, mainGroup)
	routes.Permissions(app, mainGroup)
	return router
}
