package handler

import (
	"api/internal/sensors_data/domain"
	"api/internal/sensors_data/usecase"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/tentone/mssql-uuid"
)

// Interface for handling HTTP requests related to sensor's data
type SensorDataHandler interface {
	AddSensorData(c *gin.Context)
}

// Structure request for sensor data
type SensorDataRequest struct {
	// Uuid of the sensor that recorded the data
	SensorUuid uuid.UUID `json:"sensorUuid"`
	// Array of timestamp-value pairs
	Readings []struct {
		// ISO 8601 format
		Timestamp string  `json:"timestamp"`
		Value     float64 `json:"value"`
	} `json:"readings"`
}

// Process HTTP requests and interaction with the SensorDataService
type SensorDataHandlerImpl struct {
	Service usecase.SensorDataService
}

func NewSensorDataHandler(service usecase.SensorDataService) SensorDataHandler {
	return &SensorDataHandlerImpl{Service: service}
}

func (h *SensorDataHandlerImpl) AddSensorData(c *gin.Context) {
	var req SensorDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Readings) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one sensor data is required"})
		return
	}

	var sensorDataList []*domain.SensorData
	var responseData [][]float64

	for _, reading := range req.Readings {
		// Convert timestamp string to time.Time
		timestamp, err := time.Parse(time.RFC3339, reading.Timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid timestamp format: %s", reading.Timestamp)})
			return
		}

		sensorData := &domain.SensorData{
			SensorUuid: req.SensorUuid,
			Timestamp:  timestamp,
			Value:      req.Readings[1].Value,
		}

		sensorDataList = append(sensorDataList, sensorData)
		responseData = append(responseData, []float64{float64(timestamp.Unix()), reading.Value})
	}

	var err = h.Service.AddSensorData(c.Request.Context(), sensorDataList)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid": req.SensorUuid,
		"data": responseData,
	})
}
