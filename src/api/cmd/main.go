package main

import (
	"log"

	routes_authentication "api/internal/auth"
	routes_sensors "api/internal/sensors"
	routes_sensors_data "api/internal/sensors_data"
	routes_users "api/internal/users"
	"api/version"

	_ "api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	// Middleware that allows CORS and custom headers (e.g., Authorization)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Role"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Authorization"},
	}))

	// All routes
	routes_sensors.RegisterSensorRoutes(router)
	routes_sensors_data.RegisterSensordataRoutes(router)
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
