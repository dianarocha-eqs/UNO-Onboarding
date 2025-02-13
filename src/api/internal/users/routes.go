package routes

import (
	auth_repository "api/internal/auth/repository"
	auth_service "api/internal/auth/usecase"
	users_handler "api/internal/users/handler"
	users_repository "api/internal/users/repository"
	users_service "api/internal/users/usecase"
	"api/utils"

	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Declares the routes that can be accessed for users management.
func RegisterUsersRoutes(router *gin.Engine) {

	authRepo, err := auth_repository.NewAuthRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	usersRepos, err := users_repository.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	userService := users_service.NewUserService(usersRepos)
	authService := auth_service.NewAuthService(authRepo, usersRepos)

	h := users_handler.NewUserHandler(authService, userService)

	router.Use(cors.Default())
	// User routes
	api := router.Group("/v1/users/")
	api.Use(utils.AuthMiddleware(authService))
	{
		// Create User (if admin)
		api.POST("create", utils.AdminOnly(), h.AddUser)
		// Edit User (if admin or user themself)
		api.POST("edit", utils.AdminAndUserItself(), h.EditUser)
		// List Users
		api.POST("list", h.ListUsers)
	}

	recover := router.Group("/v1/users/")
	// No authentication required
	// Reset Password
	recover.POST("change-password", h.ResetPassword)

}
