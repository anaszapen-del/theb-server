package handlers

import (
	"net/http"

	"theb-backend/internal/logger"
	"theb-backend/internal/service/auth/dtos"
	"theb-backend/internal/service/auth/services"
	apperrors "theb-backend/pkg/errors"
	"theb-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// PassengerSignup handles passenger signup request
// @Summary Passenger signup with phone number
// @ID auth-passenger-signup
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.SignupRequest true "Signup request with name and phone number"
// @Success 201 {object} dtos.SignupResponse "OTP sent successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "Phone number already exists"
// @Failure 429 {object} response.ErrorResponse "Rate limit exceeded"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Failure 503 {object} response.ErrorResponse "SMS service unavailable"
// @Router /api/v1/auth/passenger/signup [post]
func (h *AuthHandler) PassengerSignup(c *gin.Context) {
	var req dtos.SignupRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid signup request", map[string]interface{}{
			"error": err.Error(),
		})
		response.BadRequest(c, "Invalid request body", map[string]interface{}{
			"validation_error": err.Error(),
		})
		return
	}

	// Store name in context for verification step
	c.Set("signup_name", req.Name)

	// Call service
	resp, err := h.authService.PassengerSignup(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, resp, "")
}

// PassengerVerify handles OTP verification for passenger signup
// @Summary Verify OTP and complete passenger signup
// @ID auth-passenger-verify
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dtos.VerifyRequest true "Verification request with phone number and OTP"
// @Success 200 {object} dtos.VerifyResponse "Signup completed, user logged in"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired OTP"
// @Failure 409 {object} response.ErrorResponse "User already exists"
// @Failure 429 {object} response.ErrorResponse "Too many attempts"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /api/v1/auth/passenger/verify [post]
func (h *AuthHandler) PassengerVerify(c *gin.Context) {
	var req dtos.VerifyRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Invalid verify request", map[string]interface{}{
			"error": err.Error(),
		})
		response.BadRequest(c, "Invalid request body", map[string]interface{}{
			"validation_error": err.Error(),
		})
		return
	}

	// Call service with name from request
	resp, err := h.authService.VerifyPassengerSignup(c.Request.Context(), &req, req.Name)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, resp, "Signup completed successfully")
}

// handleError handles application errors and sends appropriate responses
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message, nil)
		return
	}

	logger.Error("Unhandled error in auth handler", map[string]interface{}{
		"error": err.Error(),
	})
	response.InternalServerError(c, "An unexpected error occurred")
}
