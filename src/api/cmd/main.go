package main

import (
	"api/internal/sensors/handler"
	"api/internal/sensors/repository"
	"api/internal/sensors/usecase"
	"log"

	routes "api/internal/sensors"

	"github.com/gin-gonic/gin"
)

func main() {

	var repo repository.SensorRepository
	var service usecase.SensorService
	var h handler.SensorHandler
	var err error

	repo, err = repository.NewSensorRepository()
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	service = usecase.NewSensorService(repo)

	h = handler.NewSensorHandler(service)

	router := gin.Default()
	// Sensor routes
	routes.RegisterSensorRoutes(router, h)

	log.Println("Server is running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
