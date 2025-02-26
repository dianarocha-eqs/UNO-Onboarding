package main

import (
	"log"

	routes_authentication "api/internal/auth"
	routes_sensors "api/internal/sensors"
	routes_users "api/internal/users"
	"api/version"

	_ "api/docs"

	"github.com/swaggo/files" // This replaces `swaggerFiles`
	"github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @Title Uno-Onboarding api in go
// @Version 1.0.0
//
// @Description Declares all routes from the api
//
// @Tag Users
// @Tag Auth
// @Tag Sensor
// @Tag SensorData
// @host localhost:8080
func main() {

	router := gin.Default()
	router.Use(cors.Default())

	// Sensor routes
	routes_sensors.RegisterSensorRoutes(router)
	routes_users.RegisterUsersRoutes(router)
	routes_authentication.RegisterAuthRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get api version
	router.GET("/version", version.GetVersion)

	log.Println("Server is running on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
