package dtos

import "github.com/google/uuid"

// SignupRequest represents the passenger signup request
// @name SignupRequest
type SignupRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// SignupResponse represents the signup response after OTP is sent
// @name SignupResponse
type SignupResponse struct {
	Message     string `json:"message"`
	PhoneNumber string `json:"phone_number"`
	ExpiresIn   int    `json:"expires_in"` // in seconds
}

// VerifyRequest represents the OTP verification request
// @name VerifyRequest
type VerifyRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OTPCode     string `json:"otp_code" binding:"required,len=6"`
	Name        string `json:"name" binding:"required,min=2,max=100"` // Name for signup flow
}

// VerifyResponse represents the verification response with tokens
// @name VerifyResponse
type VerifyResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	UserID       uuid.UUID `json:"user_id"`
	PhoneNumber  string    `json:"phone_number"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	ExpiresIn    int       `json:"expires_in"` // in seconds
}
