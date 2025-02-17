package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/tentone/mssql-uuid"
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
	Name    string    `json:"name"`
	UUID    uuid.UUID `json:"id"`
	Picture string    `json:"picture"`
}

// Structure request for list users
type FilterSearchAndSort struct {
	Search string `json:"search"`
	Sort   int    `json:"sort"`
}

// Process HTTP requests and interaction with the UserService for user operations
type UserHandlerImpl struct {
	Service usecase.UserService
}

func NewUserHandler(service usecase.UserService) UserHandler {
	return &UserHandlerImpl{Service: service}
}

func (h *UserHandlerImpl) AddUser(c *gin.Context) {

	// Gets token from header
	tokenAuth, _ := c.Get("token")

	// Gets role from header
	roleAuth, _ := c.Get("role")

	str := tokenAuth.(string)
	var role bool
	// Checks if the role from header is the same as the role given to the user
	err := h.Service.GetRoutesAuthorization(c.Request.Context(), str, &role, nil)
	if err != nil || role != roleAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current user is not authorized to create a new user"})
		return
	}

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
	c.JSON(http.StatusCreated, gin.H{"uuid": ID})
}

func (h *UserHandlerImpl) EditUser(c *gin.Context) {

	// Gets token from header
	tokenStr, _ := c.Get("token")

	// Gets role from header
	roleAuth, _ := c.Get("role")

	// Gets uuid from header
	uuidAuth, _ := c.Get("uuid")

	str := tokenStr.(string)
	var role bool
	var userID uuid.UUID
	// Checks if the role or uuid from header is the same as the role and uuid given to/from the user
	err := h.Service.GetRoutesAuthorization(c.Request.Context(), str, &role, &userID)
	if err != nil || role != roleAuth || userID != uuidAuth {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current user is not authorized to edit this user"})
		return
	}

	var user domain.User
	// Bind JSON to user struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call UpdateUser service
	err = h.Service.UpdateUser(c.Request.Context(), &user)
	if err != nil {
		fmt.Println(user)
		// this looks weird but i don't know how different should it be
		if strings.Contains(err.Error(), "name, email, and phone") {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusOK)
}

func (h *UserHandlerImpl) ListUsers(c *gin.Context) {

	var filter FilterSearchAndSort
	var err error
	if err = c.ShouldBindJSON(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var users []domain.User
	users, err = h.Service.ListUsers(c.Request.Context(), filter.Search, filter.Sort)
	if err != nil {
		// Check specific errors for handling them
		if strings.Contains(err.Error(), "invalid sort direction") || strings.Contains(err.Error(), "no result was found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Prepare the response
	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			Name:    user.Name,
			UUID:    user.ID,
			Picture: user.Picture,
		})
	}

	// Return the users in the expected format
	c.JSON(http.StatusOK, response)
}
