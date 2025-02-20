package usecase

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/repository"
	"context"
	"errors"
	"fmt"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's services
type SensorService interface {
	// Updates an existing sensor
	UpdateSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error
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

func (s *SensorServiceImpl) UpdateSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error {

	var stateOwner, err = s.Repo.GetSensorOwner(ctx, sensor.ID, userUuid)
	if err != nil && !stateOwner {
		return fmt.Errorf("this sensor is not from this user")
	}

	if err = validateRequiredFields(sensor); err != nil {
		return err
	}

	validColors := map[string]bool{
		domain.RED:    true,
		domain.GREEN:  true,
		domain.BLUE:   true,
		domain.YELLOW: true,
	}

	if !validColors[sensor.Color] {
		return errors.New("cannot change color if not for one of this: must be RED, GREEN, BLUE, or YELLOW")
	}

	err = s.Repo.UpdateSensor(ctx, sensor)
	if err != nil {
		return fmt.Errorf("failed to update sensor on database")
	}

	return nil
}
