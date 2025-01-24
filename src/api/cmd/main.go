package main

import (
	"log"

	routes_s "api/internal/sensors"
	"api/internal/sensors/handler"
	"api/internal/sensors/repository"
	"api/internal/sensors/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	var repos repository.SensorRepository
	var service usecase.SensorService
	var h handler.SensorHandler
	var err error

	repos, err = repository.NewSensorRepository()

	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	service = usecase.NewSensorService(repos)

	h = handler.NewSensorHandler(service)

	// Sensor routes
	routes_s.RegisterSensorRoutes(router, h)

	log.Println("Server is running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
