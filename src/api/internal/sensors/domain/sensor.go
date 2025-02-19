package domain

import (
	"api/internal/users/domain"

	uuid "github.com/tentone/mssql-uuid"
)

const (
	// Visibility constants
	PUBLIC  bool = false
	PRIVATE bool = true

	// Color constants
	RED    string = "#FF0000"
	YELLOW string = "#FFFF00"
	GREEN  string = "#00FF00"
	BLUE   string = "#OOOOFF"

	// Category constants
	TEMPERATURE int = 0
	HUMIDITY    int = 1
	PRESSURE    int = 2
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
	// Visibility: public (false) or private (true)
	Visibility bool `json:"visibility" gorm:"column:visibility;type:bit;default:1"`
	// UUID of the user who owns the sensor
	SensorOwner uuid.UUID `json:"sensor_owner" gorm:"column:sensor_owner;type:uniqueidentifier;not null"`
	// This field establishes a relationship between Sensor and User using the foreign key
	User domain.User `json:"user" gorm:"foreignKey:SensorOwner;references:ID;constraint:OnDelete:CASCADE"`
}
