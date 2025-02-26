package users

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

// @Summary Register user management routes
// @Description Declares the routes that can be accessed for user management
// @Tags users
// @Router /v1/users [post]
func RegisterUsersRoutes(router *gin.Engine) {

	authRepo, err := auth_repository.NewAuthRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	usersRepos, err := users_repository.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	userService := users_service.NewUserService(usersRepos, authRepo)
	authService := auth_service.NewAuthService(authRepo, usersRepos)

	h := users_handler.NewUserHandler(authService, userService)

	router.Use(cors.Default())
	// User routes
	api := router.Group("/v1/users/")
	api.Use(utils.AuthMiddleware(authService))
	{
		// @Summary Route to create new user
		// @Description Route to create new user (if admin only)
		// @Router /v1/users/create [post]
		api.POST("create", utils.AdminOnly(), h.AddUser)

		// @Summary Route to edit user
		// @Description Route to edit user details (if admin or user themselves)
		// @Router /v1/users/edit [post]
		api.POST("edit", utils.AdminAndUserItself(), h.EditUser)

		// @Summary Route to list users
		// @Description Route to list all users with optional filtering and sorting
		// @Router /v1/users/list [post]
		api.POST("list", h.ListUsers)
	}

	recover := router.Group("/v1/users/")
	// No authentication required
	// @Summary Route to recover password
	// @Description Route to recover password if forgotten (only receives the user's email)
	// @Router /v1/users/forgot-password [post]
	recover.POST("forgot-password", h.RecoverPassword)
	// @Summary Route to reset password
	// @Description Route to reset previous password and update with a new one
	// @Router /v1/users/change-password [post]
	recover.POST("change-password", h.ResetPassword)

}
