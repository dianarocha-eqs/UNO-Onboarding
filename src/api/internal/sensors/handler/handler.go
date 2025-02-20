package handler

import (
	sensor_service "api/internal/sensors/usecase"
	user_service "api/internal/users/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Interface for handling HTTP requests related to sensors
type SensorHandler interface {
	// GetSensors handles the retrieval of all sensors.
	ListSensors(c *gin.Context)
}

// Structure request for list sensors
type FilterSearch struct {
	// Search term to filter sensors by name
	Search string `json:"search"`
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
func (h *SensorHandlerImpl) ListSensors(c *gin.Context) {

	// Gets token from header
	var tokenAuth, _ = c.Get("token")

	var str = tokenAuth.(string)

	// Get user id from token (set by login)
	var userUuid, err = h.UserService.GetUserByToken(c.Request.Context(), str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	var req FilterSearch
	sensors, err := h.SensorService.ListSensors(c.Request.Context(), userUuid, req.Search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to list sensors"})
		return
	}
	c.JSON(http.StatusOK, sensors)
}
