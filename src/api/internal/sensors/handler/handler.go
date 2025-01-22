package handler

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/usecase"
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
	// AddSensor handles the creation of a new sensor.
	AddSensor(c *gin.Context)
	// UpdateSensor handles updating an existing sensor's details.
	UpdateSensor(c *gin.Context)
	// DeleteSensor handles deleting a sensor by its ID.
	DeleteSensor(c *gin.Context)
}

// SensorHandler handles HTTP requests related to sensors.
type SensorHandlerImpl struct {
	Service usecase.SensorService
}

func NewSensorHandler(service usecase.SensorService) SensorHandler {
	return &SensorHandlerImpl{Service: service}
}

func (h *SensorHandlerImpl) GetSensors(c *gin.Context) {
	sensors, err := h.Service.GetAllSensors()
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
	sensor, err := h.Service.GetSensorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
		return
	}
	c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandlerImpl) AddSensor(c *gin.Context) {
	var sensor domain.Sensor
	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.CreateSensor(&sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create sensor"})
		return
	}

	c.JSON(http.StatusCreated, sensor)
}

func (h *SensorHandlerImpl) UpdateSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	existingSensor, err := h.Service.GetSensorByID(uint(id))
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

	err = h.Service.UpdateSensor(&existingSensor)
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

	err = h.Service.DeleteSensor(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete sensor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted"})
}
