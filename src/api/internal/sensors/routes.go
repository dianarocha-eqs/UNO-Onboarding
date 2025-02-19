package routes

import (
	"api/internal/sensors/handler"
	sensor_repository "api/internal/sensors/repository"
	sensor_service "api/internal/sensors/usecase"
	users_repository "api/internal/users/repository"
	users_service "api/internal/users/usecase"
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

	userService := users_service.NewUserService(usersRepos)
	sensorService := sensor_service.NewSensorService(sensorRepo)

	h := handler.NewSensorHandler(sensorService, userService)

	// Sensor routes
	api := router.Group("/v1/sensor")
	{
		api.GET("sensors", h.GetSensors)          // Get all sensors
		api.GET("sensors/:id", h.GetSensor)       // Get sensor by ID
		api.POST("create", h.CreateSensor)        // Add sensor
		api.PUT("sensors/:id", h.UpdateSensor)    // Update sensor
		api.DELETE("sensors/:id", h.DeleteSensor) // Delete sensor
	}
}
