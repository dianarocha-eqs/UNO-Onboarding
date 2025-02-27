package sensors_data

import (
	auth_repository "api/internal/auth/repository"
	auth_service "api/internal/auth/usecase"
	"api/internal/sensors_data/handler"
	sensor_repository "api/internal/sensors_data/repository"
	sensor_service "api/internal/sensors_data/usecase"
	user_repository "api/internal/users/repository"
	user_service "api/internal/users/usecase"
	middleware "api/utils"

	"log"

	"github.com/gin-gonic/gin"
)

// RegisterSensorRoutes declares the routes that can be accessed for sensor management.
func RegisterSensordataRoutes(router *gin.Engine) {

	sensorDataRepo, err := sensor_repository.NewSensorDataRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	usersRepos, err := user_repository.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	authRepo, err := auth_repository.NewAuthRepository()
	if err != nil {
		log.Fatalf("Failed to create auth repository: %v", err)
	}

	_ = user_service.NewUserService(usersRepos, authRepo)
	sensorDataService := sensor_service.NewSensorDataService(sensorDataRepo)
	authService := auth_service.NewAuthService(authRepo, usersRepos)

	h := handler.NewSensorDataHandler(sensorDataService)

	// Sensor's data routes
	api := router.Group("/v1/sensor/data/")
	api.Use(middleware.AuthMiddleware(authService))
	{
		// Add sensor's data
		api.POST("add", h.AddSensorData)
		// Read sensor data
		api.POST("get", h.ReadSensorData)
	}
}
