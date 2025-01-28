package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Interface for handling HTTP requests related to users
type UserHandler interface {
	// Handles the HTTP request to create a new user
	AddUser(c *gin.Context)
	EditUser(c *gin.Context)
}

// Process HTTP requests and interaction with the UserService for user operations
type UserHandlerImpl struct {
	Service usecase.UserService
}

func NewUserHandler(service usecase.UserService) UserHandler {
	return &UserHandlerImpl{Service: service}
}

func (h *UserHandlerImpl) AddUser(c *gin.Context) {

	// Check data received on the JSON body
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the user
	ID, err := h.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		// Check if it's a validation error (missing fields)
		if strings.Contains(err.Error(), "name, email, and phone are required fields") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// Internal error (in case of duplicate emailm p.ex)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		}
		return
	}

	// Respond with the created user's uuid
	c.JSON(http.StatusCreated, gin.H{"userId": ID})
}

func (h *UserHandlerImpl) EditUser(c *gin.Context) {
	var user domain.User

	// Check data received on the JSON body
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Proceed with the user update
	if err := h.Service.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user", "error": err.Error()})
		return
	}

	// Respond with the updated user information
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
