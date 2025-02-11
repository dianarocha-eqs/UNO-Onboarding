package domain

import "encoding/json"

// Sensor represents a device that collects and transmits data about its environment.
type Sensor struct {
	// ID is the unique identifier for the sensor.
	ID uint `json:"id"`
	// Name is the human-readable name assigned to the sensor.
	Name string `json:"name"`
	// Category specifies the type of data the sensor collects, such as Temperature, Humidity, or Pressure.
	Category string `json:"category"`
	// Color assigned to the sensor.
	Color string `json:"color"`
	// Description provides additional information about the sensor's purpose or functionality.
	Description string `json:"description"`
	// Visibility defines whether the sensor is publicly visible or private.
	Visibility bool `json:"visibility"`
}

func (s *Sensor) MarshalJSON() ([]byte, error) {
	type Alias Sensor
	return json.Marshal(&struct {
		Visibility int `json:"visibility"`
		*Alias
	}{
		Visibility: boolToInt(s.Visibility),
		Alias:      (*Alias)(s),
	})
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
