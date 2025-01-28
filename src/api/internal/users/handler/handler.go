package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Interface for handling HTTP requests related to users
type UserHandler interface {
	// Handles the HTTP request to create a new user
	AddUser(c *gin.Context)
	// Handles the HTTP request to edit the info from a user
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

	// Set the default role if not provided
	if user.Role == false {
		user.Role = false
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

	// Bind the JSON body to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// These fields are required and should not be empty
	if user.Name == "" || user.Phone == "" || user.Email == "" {
		fmt.Fprintln(c.Writer, "name, phone and email cannot be empty. Previous values remained.")
	}

	// Validate the password if it is provided in the body (it should not be empty if provided)
	if user.Password != "" && user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be empty"})
		return
	}

	// Validate UUID format
	if _, err := uuid.Parse(user.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// Validate password (if provided)
	if user.Password != "" && len(user.Password) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password cannot be empty"})
		return
	}

	// If picture is empty in the body, set it to an empty string
	if user.Picture == "" {
		user.Picture = ""
	}

	// Proceed with the user update
	if err := h.Service.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user", "error": err.Error()})
		return
	}

	// Respond with the updated user information
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
