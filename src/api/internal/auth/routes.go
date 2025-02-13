package auth

import (
	"api/internal/auth/handler"
	auth_repos "api/internal/auth/repository"
	auth_service "api/internal/auth/usecase"
	user_repos "api/internal/users/repository"
	user_service "api/internal/users/usecase"
	middleware "api/utils"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up authentication routes
func RegisterAuthRoutes(router *gin.Engine) {
	// Initialize repositories
	authRepo, err := auth_repos.NewAuthRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	userRepo, err := user_repos.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	authService := auth_service.NewAuthService(authRepo, userRepo)
	userService := user_service.NewUserService(userRepo)

	h := handler.NewAuthHandler(authService, userService)

	router.Use(cors.Default())

	// Public route (no need to protect)
	auth := router.Group("/v1/auth")
	{
		auth.POST("/login", h.Login)
	}
	protect := router.Group("/v1/")
	protect.Use(middleware.AuthMiddleware(authService))
	{
		protect.POST("/auth/logout", h.Logout)
	}
}
