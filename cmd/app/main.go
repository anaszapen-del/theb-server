package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"theb-backend/internal/app"
	"theb-backend/internal/config"
	"theb-backend/internal/db"
	"theb-backend/internal/logger"
)

// @title THEB API
// @version 1.0
// @description THEB Ride-Hailing Backend API for Mafraq, Jordan
// @termsOfService http://swagger.io/terms/

// @contact.name THEB API Support
// @contact.url http://theb.app/support
// @contact.email support@theb.app

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger.Init(cfg.Logging.Level, cfg.Logging.Format)
	logger.Info("Starting THEB Backend Server", map[string]interface{}{
		"env":  cfg.App.Env,
		"port": cfg.App.Port,
	})

	// Initialize database connection
	database, err := db.InitPostgres(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	logger.Info("Database connection established", nil)

	// Run database migrations
	if err := db.AutoMigrate(database); err != nil {
		logger.Fatal("Failed to run database migrations", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Initialize Redis connection (optional in development)
	redisClient, err := db.InitRedis(cfg)
	if err != nil {
		if cfg.App.Env == "production" {
			logger.Fatal("Failed to connect to Redis", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			logger.Warn("Redis connection failed - continuing without Redis (development mode)", map[string]interface{}{
				"error": err.Error(),
			})
			redisClient = nil
		}
	} else {
		logger.Info("Redis connection established", nil)
	}

	// Initialize application
	application, err := app.New(cfg, database, redisClient)
	if err != nil {
		logger.Fatal("Failed to initialize application", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Start HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.App.Host, cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      application.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", map[string]interface{}{
			"address": addr,
		})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server shutting down...", nil)

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Close database connections
	sqlDB, _ := database.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	// Close Redis connection
	if redisClient != nil {
		redisClient.Close()
	}

	logger.Info("Server exited gracefully", nil)
}
