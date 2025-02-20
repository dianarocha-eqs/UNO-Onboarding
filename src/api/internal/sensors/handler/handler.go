package handler

import (
	"api/internal/sensors/domain"
	sensor_service "api/internal/sensors/usecase"
	user_service "api/internal/users/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SensorHandler defines the contract for handling sensor-related HTTP requests.
type SensorHandler interface {
	// GetSensors handles the retrieval of all sensors.
	GetSensors(c *gin.Context)
	// GetSensor handles the retrieval of a specific sensor by its ID.
	GetSensor(c *gin.Context)
	// Handles the HTTP request to create a new sensor
	CreateSensor(c *gin.Context)
	// UpdateSensor handles updating an existing sensor's details.
	UpdateSensor(c *gin.Context)
	// DeleteSensor handles deleting a sensor by its ID.
	DeleteSensor(c *gin.Context)
}

// SensorHandler handles HTTP requests related to sensors.
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

func (h *SensorHandlerImpl) GetSensors(c *gin.Context) {
	sensors, err := h.SensorService.GetAllSensors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve sensors"})
		return
	}
	c.JSON(http.StatusOK, sensors)
}

func (h *SensorHandlerImpl) GetSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	sensor, err := h.SensorService.GetSensorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
		return
	}
	c.JSON(http.StatusOK, sensor)
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

func (h *SensorHandlerImpl) UpdateSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	existingSensor, err := h.SensorService.GetSensorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Sensor not found"})
		return
	}

	var sensor domain.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sensor.Name != existingSensor.Name {
		existingSensor.Name = sensor.Name
	}
	if sensor.Category != existingSensor.Category {
		existingSensor.Category = sensor.Category
	}

	if sensor.Description != existingSensor.Description {
		existingSensor.Description = sensor.Description
	}

	if sensor.Visibility != existingSensor.Visibility {
		existingSensor.Visibility = sensor.Visibility
	}

	err = h.SensorService.UpdateSensor(&existingSensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update sensor"})
		return
	}

	c.JSON(http.StatusOK, existingSensor)
}

func (h *SensorHandlerImpl) DeleteSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.SensorService.DeleteSensor(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete sensor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted"})
}
