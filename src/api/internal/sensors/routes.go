package routes

import (
	auth_repository "api/internal/auth/repository"
	auth_service "api/internal/auth/usecase"
	"api/internal/sensors/handler"
	sensor_repository "api/internal/sensors/repository"
	sensor_service "api/internal/sensors/usecase"
	users_repository "api/internal/users/repository"
	users_service "api/internal/users/usecase"
	"api/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// RegisterSensorRoutes declares the routes that can be accessed for sensor management.
func RegisterSensorRoutes(router *gin.Engine) {

	sensorRepo, err := sensor_repository.NewSensorRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	usersRepos, err := users_repository.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	authRepo, err := auth_repository.NewAuthRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	_ = users_service.NewUserService(usersRepos, authRepo)
	sensorService := sensor_service.NewSensorService(sensorRepo)
	authService := auth_service.NewAuthService(authRepo, usersRepos)

	h := handler.NewSensorHandler(sensorService)

	// Sensor routes
	api := router.Group("/v1/sensor")
	api.Use(utils.AuthMiddleware(authService))
	{
		// Update sensor
		api.POST("edit", h.UpdateSensor)
	}
}
