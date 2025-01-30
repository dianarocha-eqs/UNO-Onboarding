package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"fmt"
	"net/http"
	"strings"

	uuid "github.com/tentone/mssql-uuid"

	"github.com/gin-gonic/gin"
)

// Interface for handling HTTP requests related to users
type UserHandler interface {
	// Handles the HTTP request to create a new user
	AddUser(c *gin.Context)
	// Handles the HTTP request to edit the info from a user
	EditUser(c *gin.Context)
	//  Handles the HTTP request to list users
	ListUsers(c *gin.Context)
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

	// Bind the JSON body to the user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	fmt.Println(user)

	// Validate UUID format
	if _, err := uuid.FromString(user.ID.String()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	// User update
	if err := h.Service.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}

	// Respond with the updated user information and the UUID, return 200 OK on success
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (h *UserHandlerImpl) ListUsers(c *gin.Context) {
	// Use MustGet since we know that the middleware sent those values
	search := c.MustGet("search").(string)
	sortDirection := c.MustGet("sortDirection").(int)

	// Validate the sorting value
	if sortDirection != 0 && sortDirection != 1 && sortDirection != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sortDirection. It must be 1 (ascending), -1 (descending), or 0 (no sorting)."})
		return
	}

	// Get users based on the provided filter and sorting
	users, err := h.Service.GetUsers(c.Request.Context(), search, sortDirection)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Prepare the response with UUIDs formatted as strings
	var response []map[string]string
	for _, user := range users {
		response = append(response, map[string]string{
			"name":    user.Name,
			"uuid":    user.ID.String(), // Convert UUID to string
			"picture": user.Picture,
		})
	}

	// Return the list of users
	c.JSON(http.StatusOK, response)
}
