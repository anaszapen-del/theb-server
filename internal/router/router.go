package router

import (
	"net/http"

	"theb-backend/internal/config"
	"theb-backend/internal/container"
	"theb-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// New creates and configures the router
func New(cfg *config.Config, ctn *container.Container) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS(cfg))
	router.Use(middleware.RequestID())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    cfg.App.Name,
			"env":    cfg.App.Env,
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// TODO: Register service routes here
		// Example:
		// auth := v1.Group("/auth")
		// {
		//     authHandler := getAuthHandler(ctn)
		//     auth.POST("/login", authHandler.Login)
		//     auth.POST("/verify", authHandler.Verify)
		//     auth.POST("/refresh", authHandler.Refresh)
		// }

		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}

	// WebSocket routes
	ws := router.Group("/ws")
	{
		// TODO: Register WebSocket routes here
		// ws.GET("/location/stream", wsHandler.LocationStream)
		// ws.GET("/rides/:id", wsHandler.RideUpdates)

		ws.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "websocket ready"})
		})
	}

	return router
}
