package db

import (
	"theb-backend/internal/logger"
	"theb-backend/internal/service/auth/models"

	"gorm.io/gorm"
)

// AutoMigrate runs all database migrations
func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running database migrations...", nil)

	// Migrate User model
	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.Error("Failed to migrate User model", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	logger.Info("Database migrations completed successfully", nil)
	return nil
}
