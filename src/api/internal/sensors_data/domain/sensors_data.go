package domain

import (
	"time"

	uuid "github.com/tentone/mssql-uuid"
)

// SensorData represents a recorded data point from a sensor
type SensorData struct {
	// UUID of the sensor that recorded this data
	SensorUuid uuid.UUID `json:"sensorUuid"`
	// Timestamp when the data was recorded
	Timestamp time.Time `json:"timestamp"`
	// The measured value from the sensor
	Value float64 `json:"value"`
}
