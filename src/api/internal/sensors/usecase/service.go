package usecase

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/repository"
	"context"
	"errors"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's services
type SensorService interface {
	// Creates a new sensor
	CreateSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error
}

// Handles sensor's logic and interaction with the repository
type SensorServiceImpl struct {
	Repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

// Checks the required fields of the Sensor
func validateRequiredFields(sensor *domain.Sensor) error {
	if sensor.Name == "" || (sensor.Category != domain.TEMPERATURE && sensor.Category != domain.PRESSURE && sensor.Category != domain.HUMIDITY) {
		return errors.New("name is required and category must be one of the predefined values: Temperature, Pressure or Humidity")
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

	validColors := map[string]bool{
		domain.RED:    true,
		domain.GREEN:  true,
		domain.BLUE:   true,
		domain.YELLOW: true,
	}

	// color is not required, but if selected one, it needs to be one of the predefined colors
	if sensor.Color != "" && !validColors[sensor.Color] {
		return errors.New("invalid color: must be RED, GREEN, BLUE, or YELLOW")
	}

	err = s.Repo.CreateSensor(ctx, sensor)
	if err != nil {
		return errors.New("failed to create sensor")
	}

	return nil
}
