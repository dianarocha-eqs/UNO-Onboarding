package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

	// Validate data received on the route
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

	existingUser, err := h.Service.GetUserByID(c.Request.Context(), user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	// Validate data received on the route
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Name != existingUser.Name {
		existingUser.Name = user.Name
	}
	if user.Email != existingUser.Email {
		existingUser.Email = user.Email
	}
	if user.Phone != existingUser.Phone {
		existingUser.Phone = user.Phone
	}
	if user.Picture != existingUser.Picture {
		existingUser.Picture = user.Picture
	}

	// Call the service to create the user
	err = h.Service.UpdateUser(c.Request.Context(), &existingUser)
	if err != nil {
		// Check if it's a validation error (missing fields)
		if strings.Contains(err.Error(), "name, email, and phone are required fields") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// Internal error (in case of duplicate emailm p.ex)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user", "error": err.Error()})
		}
		return c.JSON(http.StatusCreated, gin.H{"updateduser": user.ID})
	}
}
