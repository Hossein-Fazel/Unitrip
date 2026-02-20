package web

import (
	"unitrip/internal/adapter/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	engine := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}
	engine.Use(cors.New(corsConfig))
	return engine
}

func RegisterRoutes(router *gin.Engine, userHandler http.AuthHandler) {
	authGroup := router.Group("/auth")
	registerAuthRoutes(authGroup, userHandler)
}


func registerAuthRoutes(group *gin.RouterGroup, authHandler http.AuthHandler) {
	group.POST("/signup", authHandler.Signup)
}
