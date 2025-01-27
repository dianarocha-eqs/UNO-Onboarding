package routes

import (
	"api/internal/users/handler"
	"api/internal/users/repository"
	"api/internal/users/usecase"

	"api/util"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Declares the routes that can be accessed for users management.
func RegisterUsersRoutes(router *gin.Engine) {

	var repos repository.UserRepository
	var service usecase.UserService
	var h handler.UserHandler
	var err error

	repos, err = repository.NewUserRepository()

	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	service = usecase.NewUserService(repos)

	h = handler.NewUserHandler(service)

	router.Use(cors.Default())
	// User routes
	api := router.Group("/v1/users/")
	{
		// Create User (if admin)
		api.POST("createUser", util.AdminOnly(), h.AddUser)
		// To test without the admin flag
		// api.POST("createUser", h.AddUser)

	}
}
