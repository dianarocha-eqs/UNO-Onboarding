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
	// Handles the HTTP request to edit sensor
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

	var sensor domain.Sensor
	var err error
	if err = c.ShouldBindJSON(&sensor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.UpdateSensor(c.Request.Context(), &sensor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update sensor"})
		return
	}

	c.Status(http.StatusOK)
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
