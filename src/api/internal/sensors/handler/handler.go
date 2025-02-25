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
	// GetSensors handles the retrieval of all sensors.
	ListSensors(c *gin.Context)
	// Handles the HTTP request to create a new sensor
	CreateSensor(c *gin.Context)
	// Handles the HTTP request to edit sensor
	EditSensor(c *gin.Context)
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

func (h *SensorHandlerImpl) CreateSensor(c *gin.Context) {

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = h.SensorService.CreateSensor(c.Request.Context(), &sensor, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *SensorHandlerImpl) EditSensor(c *gin.Context) {

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

	err = h.SensorService.EditSensor(c.Request.Context(), &sensor, userUuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sensors, err := h.SensorService.ListSensors(c.Request.Context(), userUuid, req.Search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []gin.H
	for _, sensor := range sensors {
		response = append(response, gin.H{
			"name":       sensor.Name,
			"category":   sensor.Category,
			"visibility": sensor.Visibility,
		})
	}

	c.JSON(http.StatusOK, response)
}
