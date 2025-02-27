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
	// Handles the HTTP request to add sensor data
	AddSensorData(c *gin.Context)
	// Handles the HTTP request to read sensor data
	ReadSensorData(c *gin.Context)
}

// Structure request to add sensor data
type SensorDataRequest struct {
	// Uuid of the sensor that recorded the data
	SensorUuid uuid.UUID `json:"sensorUuid"`
	// Array of timestamp-value pairs
	Readings []struct {
		// ISO 8601 format
		Timestamp time.Time `json:"timestamp"`
		Value     float64   `json:"value"`
	} `json:"readings"`
}

// Structure request to read sensor data
type SensorDataGetRequest struct {
	// Uuid for the sensor whose data is to be retrieved
	SensorUuid uuid.UUID `json:"sensorUuid"`
	// Start date of the time range
	From time.Time `json:"from"`
	// End date of the time range
	To time.Time `json:"to"`
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
		sensorData := &domain.SensorData{
			SensorUuid: req.SensorUuid,
			Timestamp:  reading.Timestamp,
			Value:      reading.Value,
		}

		sensorDataList = append(sensorDataList, sensorData)
		responseData = append(responseData, []float64{float64(reading.Timestamp.Unix()), reading.Value})
	}

	if err := h.Service.AddSensorData(c.Request.Context(), sensorDataList); err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid": req.SensorUuid,
		"data": responseData,
	})
}

func (h *SensorDataHandlerImpl) ReadSensorData(c *gin.Context) {
	var req SensorDataGetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.From.After(req.To) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "'from' timestamp must be before 'to' timestamp"})
		return
	}

	sensorData, err := h.Service.GetSensorData(c.Request.Context(), req.SensorUuid, req.From, req.To)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responseData [][]float64
	for _, data := range sensorData {
		responseData = append(responseData, []float64{
			float64(data.Timestamp.Unix()),
			data.Value,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": responseData})
}
