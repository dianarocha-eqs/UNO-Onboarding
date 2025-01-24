package routes

import (
	"api/internal/sensors/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RegisterSensorRoutes declares the routes that can be accessed for sensor management.
// It configures the routes for handling sensors in the API.
func RegisterSensorRoutes(router *gin.Engine, sensorHandler handler.SensorHandler) {
	router.Use(cors.Default())

	// Sensor routes
	api := router.Group("/api")
	{
		api.GET("/sensors", sensorHandler.GetSensors)          // Get all sensors
		api.GET("/sensors/:id", sensorHandler.GetSensor)       // Get sensor by ID
		api.POST("/sensors", sensorHandler.AddSensor)          // Add sensor
		api.PUT("/sensors/:id", sensorHandler.UpdateSensor)    // Update sensor
		api.DELETE("/sensors/:id", sensorHandler.DeleteSensor) // Delete sensor
	}
}
