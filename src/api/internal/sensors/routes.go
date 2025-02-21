package routes

import (
	"api/internal/sensors/handler"
	"api/internal/sensors/repository"
	"api/internal/sensors/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

// RegisterSensorRoutes declares the routes that can be accessed for sensor management.
func RegisterSensorRoutes(router *gin.Engine) {

	repos, err := repository.NewSensorRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	service := usecase.NewSensorService(repos)
	h := handler.NewSensorHandler(service)

	// Sensor routes
	api := router.Group("/v1/sensor/")
	{
		api.POST("favorite", h.MarkSensorAsFavorite)
	}
}
