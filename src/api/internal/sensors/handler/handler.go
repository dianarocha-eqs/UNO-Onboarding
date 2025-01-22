package handler

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SensorHandler handles HTTP requests related to sensors.
// It acts as a bridge between the HTTP layer and the business logic layer (service.go).
type SensorHandler struct {
	Service usecase.SensorService // Service provides business logic operations for sensors.
}

func NewSensorHandler(service usecase.SensorService) *SensorHandler {
	return &SensorHandler{Service: service}
}

func (h *SensorHandler) GetSensors(c *gin.Context) {
	sensors, err := h.Service.GetAllSensors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve sensors"})
		return
	}
	c.JSON(http.StatusOK, sensors)
}

func (h *SensorHandler) GetSensor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	sensor, err := h.Service.GetSensorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
		return
	}
	c.JSON(http.StatusOK, sensor)
}

func (h *SensorHandler) AddSensor(c *gin.Context) {
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

func (h *SensorHandler) UpdateSensor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var sensor domain.Sensor

	if err := c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingSensor, err := h.Service.GetSensorByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Sensor not found"})
		return
	}

	// Update fields only if they are not empty in the request body
	if sensor.Name != "" {
		existingSensor.Name = sensor.Name
	}
	if sensor.Category != "" {
		existingSensor.Category = sensor.Category
	}
	if sensor.Description != "" {
		existingSensor.Description = sensor.Description
	}
	if sensor.Visibility != "" {
		existingSensor.Visibility = sensor.Visibility
	}

	err = h.Service.UpdateSensor(&existingSensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update sensor"})
		return
	}

	c.JSON(http.StatusOK, existingSensor)
}

func (h *SensorHandler) DeleteSensor(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.DeleteSensor(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.DeleteSensor(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete sensor"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sensor deleted"})
}
