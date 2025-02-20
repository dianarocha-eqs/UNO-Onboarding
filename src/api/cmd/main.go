package main

import (
	"log"

	routes_authentication "api/internal/auth"
	routes_sensors "api/internal/sensors"
	routes_users "api/internal/users"
	"api/version"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	// Sensor routes
	routes_sensors.RegisterSensorRoutes(router)
	routes_users.RegisterUsersRoutes(router)
	routes_authentication.RegisterAuthRoutes(router)

	// Get api version
	router.GET("/version", version.GetVersion)

	log.Println("Server is running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
