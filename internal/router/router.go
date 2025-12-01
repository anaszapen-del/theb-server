package router

import (
	"net/http"

	"theb-backend/internal/config"
	"theb-backend/internal/container"
	"theb-backend/internal/middleware"
	"theb-backend/internal/service/auth"

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
		// Register auth routes
		if err := auth.RegisterRoutes(v1, ctn); err != nil {
			panic(err)
		}

		// TODO: Register other service routes here
		// Example:
		// user.RegisterRoutes(v1, ctn)
		// location.RegisterRoutes(v1, ctn)
		// order.RegisterRoutes(v1, ctn)
		// payment.RegisterRoutes(v1, ctn)
		// rating.RegisterRoutes(v1, ctn)
		// notification.RegisterRoutes(v1, ctn)
		// captain.RegisterRoutes(v1, ctn)

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
