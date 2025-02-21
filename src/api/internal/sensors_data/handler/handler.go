package handler

import (
	"api/internal/sensors_data/usecase"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/tentone/mssql-uuid"
)

// Interface for handling HTTP requests related to sensor's data
type SensorDataHandler interface {
	// Reads all data from a sensor within a specific time interval.
	ReadSensorData(c *gin.Context)
}

// Structure request to read sensor data
type SensorDataGetRequest struct {
	// Uuid for the sensor whose data is to be retrieved
	SensorUuid uuid.UUID `json:"sensorUuid"`
	// Start date of the time range
	From string `json:"from"`
	// End date of the time range
	To string `json:"to"`
}

// Process HTTP requests and interaction with the SensorDataService
type SensorDataHandlerImpl struct {
	Service usecase.SensorDataService
}

func NewSensorDataHandler(service usecase.SensorDataService) SensorDataHandler {
	return &SensorDataHandlerImpl{Service: service}
}

func (h *SensorDataHandlerImpl) ReadSensorData(c *gin.Context) {
	var req SensorDataGetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// From ISO8601 to time.Time
	fromTime, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'from' timestamp format"})
		return
	}

	toTime, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'to' timestamp format"})
		return
	}

	sensorData, err := h.Service.GetSensorData(c.Request.Context(), req.SensorUuid, fromTime, toTime)
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
