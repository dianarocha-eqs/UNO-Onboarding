package handler

import (
	"api/internal/sensors/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's data operations
type SensorHandler interface {
	MarkSensorAsFavorite(c *gin.Context)
}

// Structure request for sensor's marked as favorites
type RequestFavoriteSensors struct {
	// Sensor UUID
	SensorUuid uuid.UUID `json:"uuid"`
	// Whether the sensor is marked as favorite or not
	Favorite bool `json:"favorite"`
}

// Performs user's data operations using database/sql to interact with the database
type SensorHandlerImpl struct {
	Service usecase.SensorService
}

func NewSensorHandler(service usecase.SensorService) SensorHandler {
	return &SensorHandlerImpl{Service: service}
}

func (h *SensorHandlerImpl) MarkSensorAsFavorite(c *gin.Context) {

	var req RequestFavoriteSensors
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.MarkSensorFavorite(c.Request.Context(), req.SensorUuid, req.Favorite)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":     req.SensorUuid,
		"favorite": req.Favorite,
	})
}
