package main

import (
	"log"

	routes_s "api/internal/sensors"
	routes_u "api/internal/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	// Sensor routes
	routes_s.RegisterSensorRoutes(router)
	routes_u.RegisterUsersRoutes(router)

	log.Println("Server is running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
