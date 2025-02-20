package domain

import (
	"api/internal/users/domain"

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
	ID uuid.UUID `json:"uuid" gorm:"column:id;type:uniqueidentifier"`
	// Name of the sensor
	Name string `json:"name" gorm:"column:name;type:nvarchar(100);not null"`
	// Category specifies the type of data the sensor collects
	Category int `json:"category" gorm:"column:category;type:int;not null"`
	// Color for the sensor, stored as a hex value (e.g., "#FF00FF")
	Color string `json:"color" gorm:"column:color;type:nvarchar(7);not null"`
	// Additional information about the sensor's functionality
	Description string `json:"description" gorm:"column:description;type:nvarchar(255);"`
	// Visibility: public (true) or private (false)
	Visibility bool `json:"visibility" gorm:"column:visibility;type:bit;default:1"`
	// UUID of the user who owns the sensor
	SensorOwnerUuid uuid.UUID `json:"sensorOwnerUuid" gorm:"column:sensorOwnerUuid;type:uniqueidentifier;not null"`
	// This field establishes a relationship between Sensor and User using the foreign key
	User domain.User `json:"user" gorm:"foreignKey:SensorOwnerUuid;references:ID;constraint:OnDelete:CASCADE"`
}
