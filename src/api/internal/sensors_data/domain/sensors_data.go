package domain

import (
	"api/internal/sensors/domain"
	"time"

	uuid "github.com/tentone/mssql-uuid"
)

// SensorData represents a recorded data point from a sensor
type SensorData struct {
	// Unique identifier for the data entry
	ID uuid.UUID `json:"id" gorm:"column:id;uniqueidentifier;not null"`
	// UUID of the sensor that recorded this data
	SensorUuid uuid.UUID `json:"sensorUuid" gorm:"column:sensorUuid;type:uniqueidentifier;not null"`
	// This field establishes a relationship between Sensor and Sensor Data using the foreign key
	Sensor domain.Sensor `json:"sensor" gorm:"foreignKey:SensorUuid;references:ID;constraint:OnDelete:CASCADE"`
	// Timestamp when the data was recorded
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp;type:datetime;not null;default:CURRENT_TIMESTAMP"`
	// The measured value from the sensor
	Value float64 `json:"value" gorm:"column:value;type:float;not null"`
}
