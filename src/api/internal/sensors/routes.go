package routes

import (
	"api/internal/sensors/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterSensorRoutes(router *gin.Engine, sensorHandler *handler.SensorHandler) {
	router.Use(cors.Default())

	// Sensor routes
	api := router.Group("/api")
	{
		api.GET("/sensors", sensorHandler.GetSensors)
		api.GET("/sensors/:id", sensorHandler.GetSensor)
		api.POST("/sensors", sensorHandler.AddSensor)
		api.PUT("/sensors/:id", sensorHandler.UpdateSensor)
		api.DELETE("/sensors/:id", sensorHandler.DeleteSensor)
	}
}
