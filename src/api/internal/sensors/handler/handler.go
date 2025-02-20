package handler

import (
	"api/internal/sensors/domain"
	sensor_service "api/internal/sensors/usecase"
	user_service "api/internal/users/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface for handling HTTP requests related to sensors
type SensorHandler interface {
	// Handles the HTTP request to edit sensor
	UpdateSensor(c *gin.Context)
}

// Process HTTP requests and interaction with SensorService/UserService for sensor operations
type SensorHandlerImpl struct {
	SensorService sensor_service.SensorService
	UserService   user_service.UserService
}

func NewSensorHandler(sensorService sensor_service.SensorService, userService user_service.UserService) SensorHandler {
	return &SensorHandlerImpl{
		SensorService: sensorService,
		UserService:   userService,
	}
}

func (h *SensorHandlerImpl) UpdateSensor(c *gin.Context) {

	// Gets token from header
	var tokenAuth, _ = c.Get("token")

	var str = tokenAuth.(string)

	// Get user id from token (set by login)
	var userUuid, err = h.UserService.GetUserByToken(c.Request.Context(), str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	var sensor domain.Sensor
	if err = c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.SensorService.UpdateSensor(c.Request.Context(), &sensor, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update sensor"})
		return
	}

	c.Status(http.StatusOK)
}
