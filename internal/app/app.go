package app

import (
	"theb-backend/internal/config"
	"theb-backend/internal/container"
	"theb-backend/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Application represents the application
type Application struct {
	config    *config.Config
	container *container.Container
	router    *gin.Engine
}

// New creates a new application instance
func New(cfg *config.Config, db *gorm.DB, redis *redis.Client) (*Application, error) {
	// Set Gin mode
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize dependency injection container
	ctn := container.New()

	// Register core dependencies
	ctn.Register("config", cfg)
	ctn.Register("db", db)
	ctn.Register("redis", redis)

	// Register all services
	if err := registerServices(ctn); err != nil {
		return nil, err
	}

	// Initialize router
	r := router.New(cfg, ctn)

	return &Application{
		config:    cfg,
		container: ctn,
		router:    r,
	}, nil
}

// Router returns the Gin router
func (a *Application) Router() *gin.Engine {
	return a.router
}

// registerServices registers all service dependencies
func registerServices(ctn *container.Container) error {
	// TODO: Register service modules here
	// Example:
	// auth.RegisterService(ctn)
	// user.RegisterService(ctn)
	// location.RegisterService(ctn)
	// order.RegisterService(ctn)
	// payment.RegisterService(ctn)
	// rating.RegisterService(ctn)
	// notification.RegisterService(ctn)
	// captain.RegisterService(ctn)

	return nil
}
