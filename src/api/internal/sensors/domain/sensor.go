package domain

// Sensor represents a device that collects and transmits data about its environment.
// Each sensor has attributes like ID, name, color, category, description, and visibility.
type Sensor struct {
	ID          uint   `json:"id"`          // ID is the unique identifier for the sensor.
	Name        string `json:"name"`        // Name is the human-readable name assigned to the sensor.
	Category    string `json:"category"`    // Category specifies the type of data the sensor collects, such as Temperature, Humidity, or Pressure.
	Description string `json:"description"` // Description provides additional information about the sensor's purpose or functionality.
	Visibility  string `json:"visibility"`  // Visibility defines whether the sensor is publicly visible or private.
}
