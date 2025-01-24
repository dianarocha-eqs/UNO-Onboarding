package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	AddUser(c *gin.Context)
	// GetUsers(c *gin.Context)
}

type UserHandlerImpl struct {
	Service usecase.UserService
}

func NewUserHandler(service usecase.UserService) UserHandler {
	return &UserHandlerImpl{Service: service}
}

// func AdminOnly() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Example: Extract user role from token or context
// 		role := c.GetBool("role") // Assume this is set elsewhere
// 		if !role {                // Check if the user is an admin
// 			c.JSON(http.StatusForbidden, gin.H{"error": "403: forbidden"})
// 			c.Abort()
// 			return
// 		}
// 		c.Set("isAdmin", true)
// 		c.Next()
// 	}
// }

func (h *UserHandlerImpl) AddUser(c *gin.Context) {
	// Extract the role from the request (e.g., middleware, token, or context)
	// isAdmin := c.GetBool("isAdmin") // Assume middleware sets this value in the context
	// if !isAdmin {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "403: only admins can create users"})
	// 	return
	// }

	// Parse the incoming JSON request
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the service to create the user
	ID, err := h.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": err.Error()})
		return
	}

	// Respond with the created user's ID
	c.JSON(http.StatusCreated, gin.H{"userId": ID})
}

// func (h *UserHandlerImpl) GetUsers(c *gin.Context) {

// 	users, err := h.Service.GetAllUsers()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve sensors"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, users)
// }
