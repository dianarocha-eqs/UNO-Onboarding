package handler

import (
	auth_service "api/internal/auth/usecase"
	user_service "api/internal/users/usecase"
	"api/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface for handling HTTP requests related to authentication
type AuthHandler interface {
	// Handles the HTTP request to login a user
	Login(c *gin.Context)
	// Handles the HTTP request to logout a user
	Logout(c *gin.Context)
}

// Structure request for login
type loginRequest struct {
	// Email of the user
	Email string `json:"email" binding:"required"`
	// Password of the user
	Password string `json:"password" binding:"required"`
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
		fmt.Print(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT token
	tokenStr, err := h.AuthService.AddToken(c.Request.Context(), user)
	if err != nil {
		fmt.Print(err)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to retrieve token from authorization header"})
		return
	}

	// Call service to invalidate the token
	err := h.AuthService.InvalidateToken(c.Request.Context(), tokenStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
