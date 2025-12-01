package repositories

import (
	"context"
	"errors"
	"fmt"

	"theb-backend/internal/service/auth/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles user data access operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CheckPhoneExists checks if a phone number already exists
func (r *UserRepository) CheckPhoneExists(ctx context.Context, phoneNumber string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("phone_number = ?", phoneNumber).
		Count(&count).
		Error

	if err != nil {
		return false, fmt.Errorf("failed to check phone existence: %w", err)
	}

	return count > 0, nil
}

// CreatePassenger creates a new passenger user
func (r *UserRepository) CreatePassenger(ctx context.Context, name, phoneNumber string) (*models.User, error) {
	user := &models.User{
		ID:          uuid.New(),
		Name:        name,
		PhoneNumber: phoneNumber,
		Role:        models.RolePassenger,
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create passenger: %w", err)
	}

	return user, nil
}

// GetUserByPhone retrieves a user by phone number
func (r *UserRepository) GetUserByPhone(ctx context.Context, phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("phone_number = ?", phoneNumber).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}
