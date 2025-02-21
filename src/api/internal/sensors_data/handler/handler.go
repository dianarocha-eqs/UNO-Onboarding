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

type SensorDataHandler interface {
	AddSensorData(c *gin.Context)
}

// Structure request for sensor data
type SensorDataRequest struct {
	// sensor uuid to associate the data
	SensorUuid uuid.UUID `json:"sensorUuid"`
	// timestamp (x)
	Timestamp string `json:"timestamp"` // Timestamp in ISO 8601 format (string)
	// value (y)
	Value float64 `json:"value"`
}

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

	// Convert the timestamp from ISO 8601 format (string) to time.Time
	timestamp, err := time.Parse(time.RFC3339, req.Timestamp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timestamp format, should be ISO 8601"})
		return
	}

	sensorData := &domain.SensorData{
		ID:         uuid.NewV4(),
		SensorUuid: req.SensorUuid,
		Timestamp:  timestamp,
		Value:      req.Value,
	}

	err = h.Service.AddSensorData(c.Request.Context(), sensorData)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare the response format
	responseData := [][]float64{
		{float64(sensorData.Timestamp.Unix()), sensorData.Value},
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid": sensorData.ID,
		"data": responseData,
	})
}
