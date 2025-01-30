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

// Structure response for list users
type UserResponse struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	Picture string `json:"picture"`
}

// Process HTTP requests and interaction with the UserService for user operations
type UserHandlerImpl struct {
	Service usecase.UserService
}

func NewUserHandler(service usecase.UserService) UserHandler {
	return &UserHandlerImpl{Service: service}
}

func (h *UserHandlerImpl) AddUser(c *gin.Context) {

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID, err := h.Service.CreateUser(c.Request.Context(), &user)
	if err != nil {
		// Check if it's a validation error (missing fields)
		if strings.Contains(err.Error(), "required fields") || strings.Contains(err.Error(), "invalid email format") || strings.Contains(err.Error(), "invalid phone number format") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// Internal error (in case of duplicate email p.ex)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userId": ID})
}

func (h *UserHandlerImpl) EditUser(c *gin.Context) {
	var user domain.User

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

	if err := h.Service.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (h *UserHandlerImpl) ListUsers(c *gin.Context) {
	// Take the values from header to pass on service
	search := c.MustGet("search").(string)
	sortDirection := c.MustGet("sortDirection").(int)

	var err error
	var users []domain.User
	users, err = h.Service.GetUsers(c.Request.Context(), search, sortDirection)
	if err != nil {
		if err.Error() == "nothing with that value" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			Name:    user.Name,
			UUID:    user.ID.String(),
			Picture: user.Picture,
		})
	}
	c.JSON(http.StatusOK, response)
}
