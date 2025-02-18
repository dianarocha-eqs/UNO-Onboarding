package usecase

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/repository"
	"context"
	"errors"

	uuid "github.com/tentone/mssql-uuid"
)

// SensorService defines the business logic methods for managing sensors.
type SensorService interface {
	// CreateSensor creates a new sensor
	CreateSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error
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
	Repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

// Checks the required fields of the Sensor.
func validateRequiredFields(sensor *domain.Sensor) error {
	if sensor.Name == "" || (sensor.Category != domain.TEMPERATURE && sensor.Category != domain.PRESSURE && sensor.Category != domain.HUMIDITY) {
		return errors.New("name is required and category must be one of the predefined values (Temperature = 0, Humidity = 1, Pressure = 2 )")
	}
	return nil
}

func (s *SensorServiceImpl) CreateSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error {

	var err error
	if err = validateRequiredFields(sensor); err != nil {
		return err
	}

	sensor.ID = uuid.NewV4()
	sensor.SensorOwner = userUuid
	err = s.Repo.CreateSensor(ctx, sensor)
	if err != nil {
		return errors.New("failed to create sensor")
	}

	return s.Repo.CreateSensor(ctx, sensor)
}

func (s *SensorServiceImpl) UpdateSensor(sensor *domain.Sensor) error {
	if err := validateRequiredFields(sensor); err != nil {
		return err
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
