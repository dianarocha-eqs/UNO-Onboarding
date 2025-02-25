package domain

import (
	uuid "github.com/tentone/mssql-uuid"
)

const (
	// Color constants
	SENSOR_COLOR_RED    string = "#FF0000"
	SENSOR_COLOR_YELLOW string = "#FFFF00"
	SENSOR_COLOR_GREEN  string = "#00FF00"
	SENSOR_COLOR_BLUE   string = "#0000FF"
)

const (
	// Category constants
	SENSOR_CATEGORY_TEMPERATURE int = 0
	SENSOR_CATEGORY_HUMIDITY    int = 1
	SENSOR_CATEGORY_PRESSURE    int = 2
)

// Sensor represents a device that collects and transmits data about its environment.
type Sensor struct {
	// Unique identifier for the sensor
	ID uuid.UUID `json:"uuid"`
	// Name of the sensor
	Name string `json:"name"`
	// Category specifies the type of data the sensor collects
	Category int `json:"category"`
	// Color for the sensor
	Color string `json:"color"`
	// Additional information about the sensor's functionality
	Description string `json:"description"`
	// Visibility: public (true) or private (false)
	Visibility bool `json:"visibility"`
	// UUID of the user who owns the sensor
	SensorOwnerUuid uuid.UUID `json:"sensorOwnerUuid"`
}
