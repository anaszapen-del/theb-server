package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

// OTPService handles OTP storage and verification
type OTPService struct {
	redis  *redis.Client
	expiry time.Duration
}

// NewOTPService creates a new OTP service
func NewOTPService(redis *redis.Client, expiry time.Duration) *OTPService {
	return &OTPService{
		redis:  redis,
		expiry: expiry,
	}
}

// StoreOTP stores a hashed OTP in Redis
func (s *OTPService) StoreOTP(ctx context.Context, phoneNumber, otpCode string) error {
	// Hash the OTP code
	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otpCode), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash OTP: %w", err)
	}

	// Store in Redis with expiry
	key := fmt.Sprintf("otp:%s", phoneNumber)
	err = s.redis.Set(ctx, key, string(hashedOTP), s.expiry).Err()
	if err != nil {
		return fmt.Errorf("failed to store OTP in Redis: %w", err)
	}

	return nil
}

// VerifyOTP verifies the OTP code against stored hash
func (s *OTPService) VerifyOTP(ctx context.Context, phoneNumber, otpCode string) (bool, error) {
	key := fmt.Sprintf("otp:%s", phoneNumber)

	// Get stored hashed OTP
	hashedOTP, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("OTP not found or expired")
	}
	if err != nil {
		return false, fmt.Errorf("failed to retrieve OTP: %w", err)
	}

	// Compare OTP codes
	err = bcrypt.CompareHashAndPassword([]byte(hashedOTP), []byte(otpCode))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// DeleteOTP removes OTP from Redis
func (s *OTPService) DeleteOTP(ctx context.Context, phoneNumber string) error {
	key := fmt.Sprintf("otp:%s", phoneNumber)
	err := s.redis.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete OTP: %w", err)
	}
	return nil
}

// IncrementAttempts increments OTP verification attempts
func (s *OTPService) IncrementAttempts(ctx context.Context, phoneNumber string) (int64, error) {
	key := fmt.Sprintf("otp_attempts:%s", phoneNumber)
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment attempts: %w", err)
	}

	// Set expiry on first attempt
	if count == 1 {
		s.redis.Expire(ctx, key, time.Hour)
	}

	return count, nil
}

// GetAttempts gets the number of OTP verification attempts
func (s *OTPService) GetAttempts(ctx context.Context, phoneNumber string) (int64, error) {
	key := fmt.Sprintf("otp_attempts:%s", phoneNumber)
	count, err := s.redis.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get attempts: %w", err)
	}
	return count, nil
}
