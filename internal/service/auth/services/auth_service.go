package services

import (
	"context"
	"fmt"
	"time"

	"theb-backend/internal/config"
	"theb-backend/internal/logger"
	"theb-backend/internal/service/auth/dtos"
	"theb-backend/internal/service/auth/repositories"
	"theb-backend/internal/service/auth/utils"
	apperrors "theb-backend/pkg/errors"

	"github.com/go-redis/redis/v8"
)

// AuthService handles authentication business logic
type AuthService struct {
	config     *config.Config
	userRepo   *repositories.UserRepository
	otpService *utils.OTPService
	redis      *redis.Client
}

// NewAuthService creates a new auth service
func NewAuthService(cfg *config.Config, userRepo *repositories.UserRepository, redis *redis.Client) *AuthService {
	otpService := utils.NewOTPService(redis, cfg.OTP.Expiry)
	return &AuthService{
		config:     cfg,
		userRepo:   userRepo,
		otpService: otpService,
		redis:      redis,
	}
}

// PassengerSignup handles passenger signup with OTP
func (s *AuthService) PassengerSignup(ctx context.Context, req *dtos.SignupRequest) (*dtos.SignupResponse, error) {
	// Normalize and validate phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)
	if err := utils.ValidatePhoneNumber(phoneNumber); err != nil {
		return nil, apperrors.BadRequest("Invalid phone number format", err)
	}

	// Validate name
	if err := utils.ValidateName(req.Name); err != nil {
		return nil, apperrors.BadRequest("Invalid name", err)
	}

	// Check rate limiting for signup
	if err := s.checkSignupRateLimit(ctx, phoneNumber); err != nil {
		return nil, err
	}

	// Check if phone already exists
	exists, err := s.userRepo.CheckPhoneExists(ctx, phoneNumber)
	if err != nil {
		logger.Error("Failed to check phone existence", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to process signup", err)
	}

	if exists {
		return nil, apperrors.Conflict("Phone number already registered. Please login instead.", nil)
	}

	// Generate OTP
	otpCode, err := utils.GenerateOTP(s.config.OTP.Length)
	if err != nil {
		logger.Error("Failed to generate OTP", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to generate OTP", err)
	}

	// Store OTP in Redis
	if err := s.otpService.StoreOTP(ctx, phoneNumber, otpCode); err != nil {
		logger.Error("Failed to store OTP", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to store OTP", err)
	}

	// Send OTP via SMS (mock for now - integrate SMS provider later)
	if err := s.sendOTP(ctx, phoneNumber, otpCode); err != nil {
		logger.Error("Failed to send OTP", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
		return nil, apperrors.ServiceUnavailable("Failed to send OTP. Please try again.", nil)
	}

	logger.Info("OTP sent successfully", map[string]interface{}{
		"phone": phoneNumber,
	})

	return &dtos.SignupResponse{
		Message:     "OTP sent to phone number",
		PhoneNumber: phoneNumber,
		ExpiresIn:   int(s.config.OTP.Expiry.Seconds()),
	}, nil
}

// VerifyPassengerSignup verifies OTP and creates passenger account
func (s *AuthService) VerifyPassengerSignup(ctx context.Context, req *dtos.VerifyRequest, name string) (*dtos.VerifyResponse, error) {
	// Normalize phone number
	phoneNumber := utils.NormalizePhoneNumber(req.PhoneNumber)

	// Check verification rate limiting
	if err := s.checkVerifyRateLimit(ctx, phoneNumber); err != nil {
		return nil, err
	}

	// Verify OTP
	valid, err := s.otpService.VerifyOTP(ctx, phoneNumber, req.OTPCode)
	if err != nil {
		logger.Error("Failed to verify OTP", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})

		// Increment attempt counter
		s.otpService.IncrementAttempts(ctx, phoneNumber)

		return nil, apperrors.Unauthorized("OTP is invalid or expired", nil)
	}

	if !valid {
		// Increment attempt counter
		attempts, _ := s.otpService.IncrementAttempts(ctx, phoneNumber)
		logger.Warn("Invalid OTP attempt", map[string]interface{}{
			"phone":    phoneNumber,
			"attempts": attempts,
		})

		return nil, apperrors.Unauthorized("Invalid OTP code", nil)
	}

	// Check if user was created (shouldn't exist for signup)
	existingUser, err := s.userRepo.GetUserByPhone(ctx, phoneNumber)
	if err != nil {
		logger.Error("Failed to check existing user", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to verify user", err)
	}

	if existingUser != nil {
		// User already exists - delete OTP and return error
		s.otpService.DeleteOTP(ctx, phoneNumber)
		return nil, apperrors.Conflict("User already exists. Please login instead.", nil)
	}

	// Create passenger user
	user, err := s.userRepo.CreatePassenger(ctx, name, phoneNumber)
	if err != nil {
		logger.Error("Failed to create passenger", map[string]interface{}{
			"phone": phoneNumber,
			"name":  name,
			"error": err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to create user account", err)
	}

	// Delete OTP after successful verification
	if err := s.otpService.DeleteOTP(ctx, phoneNumber); err != nil {
		logger.Warn("Failed to delete OTP after verification", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
	}

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(
		user.ID,
		user.PhoneNumber,
		string(user.Role),
		s.config.JWT.Secret,
		s.config.JWT.AccessTokenExpiry,
	)
	if err != nil {
		logger.Error("Failed to generate access token", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(
		user.ID,
		user.PhoneNumber,
		string(user.Role),
		s.config.JWT.Secret,
		s.config.JWT.RefreshTokenExpiry,
	)
	if err != nil {
		logger.Error("Failed to generate refresh token", map[string]interface{}{
			"user_id": user.ID,
			"error":   err.Error(),
		})
		return nil, apperrors.InternalServerError("Failed to generate refresh token", err)
	}

	logger.Info("Passenger signup completed successfully", map[string]interface{}{
		"user_id": user.ID,
		"phone":   phoneNumber,
	})

	return &dtos.VerifyResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		PhoneNumber:  user.PhoneNumber,
		Name:         user.Name,
		Role:         string(user.Role),
		ExpiresIn:    int(s.config.JWT.AccessTokenExpiry.Seconds()),
	}, nil
}

// checkSignupRateLimit checks rate limit for signup requests
func (s *AuthService) checkSignupRateLimit(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("signup_attempts:%s", phoneNumber)
	count, err := s.redis.Get(ctx, key).Int64()

	if err != nil && err != redis.Nil {
		logger.Error("Failed to check signup rate limit", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
		// Don't fail request if rate limit check fails
		return nil
	}

	if count >= int64(s.config.RateLimit.OTPPerHour) {
		return apperrors.TooManyRequests("Too many signup attempts. Please try again later.", nil)
	}

	// Increment counter
	pipe := s.redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Hour)
	if _, err := pipe.Exec(ctx); err != nil {
		logger.Warn("Failed to update signup rate limit counter", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
	}

	return nil
}

// checkVerifyRateLimit checks rate limit for verification requests
func (s *AuthService) checkVerifyRateLimit(ctx context.Context, phoneNumber string) error {
	attempts, err := s.otpService.GetAttempts(ctx, phoneNumber)
	if err != nil {
		logger.Error("Failed to check verify rate limit", map[string]interface{}{
			"phone": phoneNumber,
			"error": err.Error(),
		})
		// Don't fail request if rate limit check fails
		return nil
	}

	if attempts >= 10 {
		return apperrors.TooManyRequests("Too many verification attempts. Please request a new OTP.", nil)
	}

	return nil
}

// sendOTP sends OTP via SMS provider (mock implementation)
func (s *AuthService) sendOTP(ctx context.Context, phoneNumber, otpCode string) error {
	// TODO: Integrate with SMS provider (Twilio, AWS SNS, or local Jordanian provider)
	// For development, just log the OTP
	logger.Info("SMS OTP (DEVELOPMENT MODE)", map[string]interface{}{
		"phone":    phoneNumber,
		"otp":      otpCode,
		"message":  fmt.Sprintf("Your THEB verification code is: %s", otpCode),
		"reminder": "This is logged only in development. In production, integrate SMS provider.",
	})

	return nil
}
