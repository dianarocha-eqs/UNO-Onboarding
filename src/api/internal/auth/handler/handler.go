package handler

import (
	auth_service "api/internal/auth/usecase"
	user_service "api/internal/users/usecase"
	"api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface for handling HTTP requests related to authentication
type AuthHandler interface {
	// Handles the HTTP request to login a user
	Login(c *gin.Context)
	// Handles the HTTP request to logout a user
	Logout(c *gin.Context)
	// Handles the HTTP request to reset password
	ResetPassword(c *gin.Context)
}

// Process HTTP requests and interaction with the AuthService and UserService for authentication operations
type AuthHandlerImpl struct {
	AuthService auth_service.AuthService
	UserService user_service.UserService
}

func NewAuthHandler(authService auth_service.AuthService, userService user_service.UserService) AuthHandler {
	return &AuthHandlerImpl{
		AuthService: authService,
		UserService: userService,
	}
}

// Structure request for login
type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Structure request for reset password
type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {

	var req loginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	_, hashedpassword, err := utils.GeneratePasswordHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate password hash"})
	}
	// Fetch user by email and password
	user, err := h.UserService.GetUserByEmailAndPassword(c.Request.Context(), req.Email, hashedpassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT token
	tokenStr, err := h.AuthService.AddToken(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	// Send response with user details and token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenStr,
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"role":    user.Role,
			"phone":   user.Phone,
			"picture": user.Picture,
		},
	})
}

func (h *AuthHandlerImpl) Logout(c *gin.Context) {
	// Retrieve the token from context (set by middleware)
	tokenStr, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
		return
	}

	// Call service to invalidate the token
	err := h.AuthService.InvalidateToken(c.Request.Context(), tokenStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandlerImpl) ResetPassword(c *gin.Context) {

	// Retrieve the token from context (set by middleware)
	tokenStr, exists := c.Get("token")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token or password"})
		return
	}

	if tokenStr != req.Token {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Fetch user by token
	userID, err := h.AuthService.GetUserByToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	_, hashedPassword, err := utils.GeneratePasswordHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate password hash"})
	}

	err = h.UserService.UpdatePassword(c.Request.Context(), userID, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.Status(http.StatusOK)
}
