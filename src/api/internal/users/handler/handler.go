package handler

import (
	"api/internal/users/domain"
	"api/internal/users/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	AddUser(c *gin.Context)
	GetUsers(c *gin.Context)
}

type UserHandlerImpl struct {
	Service usecase.UserService
}

func NewUserHandler(service usecase.UserService) UserHandler {
	return &UserHandlerImpl{Service: service}
}

func (h *UserHandlerImpl) AddUser(c *gin.Context) {
	var user domain.User
	var err error
	if !user.Role {
		c.JSON(http.StatusForbidden, gin.H{"error": err})
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.Service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userId": user.UserID})
}

func (h *UserHandlerImpl) GetUsers(c *gin.Context) {

	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve sensors"})
		return
	}
	c.JSON(http.StatusOK, users)
}
