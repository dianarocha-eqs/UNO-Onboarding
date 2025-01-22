package domain

type Sensor struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Visibility  string `json:"visibility"`
}
