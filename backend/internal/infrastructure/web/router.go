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

func RegisterRoutes(router *gin.Engine, userHandler http.AuthHandler, busHandler http.BusHandler, jwtService http.JWTService) {
	authGroup := router.Group("/auth")
	busGroup := router.Group("/bus")
	registerAuthRoutes(authGroup, userHandler)
	registerBusRoutes(busGroup, busHandler, jwtService)
}

func registerAuthRoutes(group *gin.RouterGroup, authHandler http.AuthHandler) {
	group.POST("/signup", authHandler.Signup)
	group.POST("/login", authHandler.Login)
}

func registerBusRoutes(group *gin.RouterGroup, busHandler http.BusHandler, jwtService http.JWTService) {
	group.GET("", busHandler.List)
	group.GET("/:id", busHandler.GetByID)
	group.POST("", AuthMiddleware(jwtService), AdminOnly(), busHandler.Create)
	group.PUT("/:id", AuthMiddleware(jwtService), AdminOnly(), busHandler.Update)
	group.DELETE("/:id", AuthMiddleware(jwtService), AdminOnly(), busHandler.Delete)
}
