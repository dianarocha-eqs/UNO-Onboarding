package handler

import (
	"api/internal/sensors_data/domain"
	"api/internal/sensors_data/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SensorDataHandler interface {
	AddSensorData(c *gin.Context)
}

type SensorDataHandlerImpl struct {
	Service usecase.SensorDataService
}

func NewSensorDataHandler(service usecase.SensorDataService) SensorDataHandler {
	return &SensorDataHandlerImpl{Service: service}
}

func (h *SensorDataHandlerImpl) AddSensorData(c *gin.Context) {
	var sensorData *domain.SensorData
	var err = h.Service.AddSensorData(c.Request.Context(), sensorData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve sensors"})
		return
	}
	c.JSON(http.StatusOK, sensorData)
}
