package usecase

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/repository"
	"context"
	"errors"
	"fmt"
)

// SensorService defines the business logic methods for managing sensors.
type SensorService interface {
	// CreateSensor creates a new sensor
	CreateSensor(sensor *domain.Sensor) error
	// DeleteSensor removes a sensor by its ID
	DeleteSensor(id uint) error
	// GetAllSensors retrieves all sensors
	GetAllSensors() ([]domain.Sensor, error)
	// GetSensorByID retrieves a sensor by its ID
	GetSensorByID(id uint) (domain.Sensor, error)
	// Updates an existing sensor
	UpdateSensor(ctx context.Context, sensor *domain.Sensor) error
}

// SensorServiceImpl is the implementation of the SensorService interface.
// It uses the SensorRepository for interacting with the underlying data storage.
type SensorServiceImpl struct {
	Repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

// Checks the required fields of the Sensor
func validateRequiredFields(sensor *domain.Sensor) error {
	if sensor.Name == "" || (sensor.Category != domain.TEMPERATURE && sensor.Category != domain.PRESSURE && sensor.Category != domain.HUMIDITY) {
		return errors.New("name is required and category must be one of the predefined values (Temperature = 0, Humidity = 1, Pressure = 2 )")
	}
	return nil
}

func (s *SensorServiceImpl) CreateSensor(sensor *domain.Sensor) error {
	if err := validateRequiredFields(sensor); err != nil {
		return err
	}
	return s.Repo.CreateSensor(sensor)
}

func (s *SensorServiceImpl) UpdateSensor(ctx context.Context, sensor *domain.Sensor) error {
	var err error
	if err = validateRequiredFields(sensor); err != nil {
		return err
	}

	err = s.Repo.UpdateSensor(ctx, sensor)
	if err != nil {
		return fmt.Errorf("failed to update sensor from id %s: %v", sensor.SensorOwner, err)
	}

	return nil
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
