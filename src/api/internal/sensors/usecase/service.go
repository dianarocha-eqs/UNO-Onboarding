package usecase

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/repository"
	"errors"
)

// SensorService defines the business logic methods for managing sensors.
// It provides operations to create, delete, retrieve, and update sensor data.
type SensorService interface {
	// CreateSensor creates a new sensor
	CreateSensor(sensor *domain.Sensor) error
	// DeleteSensor removes a sensor by its ID
	DeleteSensor(id uint) error
	// GetAllSensors retrieves all sensors
	GetAllSensors() ([]domain.Sensor, error)
	// GetSensorByID retrieves a sensor by its ID
	GetSensorByID(id uint) (domain.Sensor, error)
	// UpdateSensor updates an existing sensor
	UpdateSensor(sensor *domain.Sensor) error
}

// SensorServiceImpl is the implementation of the SensorService interface.
// It uses the SensorRepository for interacting with the underlying data storage.
type SensorServiceImpl struct {
	Repo repository.SensorRepository // Repo is the repository that handles database operations for sensors.
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

func (s *SensorServiceImpl) CreateSensor(sensor *domain.Sensor) error {
	if sensor.Name == "" && sensor.Category == "" {
		return errors.New("name and category are required")
	}
	return s.Repo.CreateSensor(sensor)
}

func (s *SensorServiceImpl) UpdateSensor(sensor *domain.Sensor) error {
	if sensor.Name == "" && sensor.Category == "" {
		return errors.New("name and category are required")
	}
	return s.Repo.UpdateSensor(sensor)
}

func (s *SensorServiceImpl) DeleteSensor(id uint) error {
	return s.Repo.DeleteSensor(id)
}

func (s *SensorServiceImpl) GetAllSensors() ([]domain.Sensor, error) {
	return s.Repo.GetAllSensors()
}

func (s *SensorServiceImpl) GetSensorByID(id uint) (domain.Sensor, error) {
	return s.Repo.GetSensorByID(id)
}
