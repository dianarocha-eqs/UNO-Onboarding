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
		// @Router /v1/users/create [post]
		api.POST("create", utils.AdminOnly(), h.AddUser)
		// @Router /v1/users/edit [post]
		api.POST("edit", utils.AdminAndUserItself(), h.EditUser)
		// @Router /v1/users/list [post]
		api.POST("list", h.ListUsers)
	}

	recover := router.Group("/v1/users/")
	// No authentication required
	// @Router /v1/users/forgot-password [post]
	recover.POST("forgot-password", h.RecoverPassword)
	// @Router /v1/users/change-password [post]
	recover.POST("change-password", h.ResetPassword)

}
