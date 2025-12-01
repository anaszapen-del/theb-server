package auth

import (
	"theb-backend/internal/config"
	"theb-backend/internal/container"
	"theb-backend/internal/service/auth/handlers"
	"theb-backend/internal/service/auth/repositories"
	"theb-backend/internal/service/auth/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// RegisterService registers auth service dependencies in the container
func RegisterService(ctn *container.Container) error {
	// Get dependencies from container
	db, err := container.GetTyped[*gorm.DB](ctn, "db")
	if err != nil {
		return err
	}

	redisClient, err := container.GetTyped[*redis.Client](ctn, "redis")
	if err != nil {
		return err
	}

	cfg, err := container.GetTyped[*config.Config](ctn, "config")
	if err != nil {
		return err
	}

	// Create instances
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(cfg, userRepo, redisClient)
	authHandler := handlers.NewAuthHandler(authService)

	// Register in container
	ctn.Register("auth.userRepo", userRepo)
	ctn.Register("auth.service", authService)
	ctn.Register("auth.handler", authHandler)

	return nil
}

// RegisterRoutes registers auth HTTP routes
func RegisterRoutes(router *gin.RouterGroup, ctn *container.Container) error {
	handler, err := container.GetTyped[*handlers.AuthHandler](ctn, "auth.handler")
	if err != nil {
		return err
	}

	auth := router.Group("/auth")
	{
		passenger := auth.Group("/passenger")
		{
			passenger.POST("/signup", handler.PassengerSignup)
			passenger.POST("/verify", handler.PassengerVerify)
		}
	}

	return nil
}
