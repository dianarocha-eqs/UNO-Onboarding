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
	// Creates a new sensor
	CreateSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error
	// Updates an existing sensor
	EditSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error
	// List all sensors
	ListSensors(ctx context.Context, userUuid uuid.UUID, search string) ([]domain.Sensor, error)
	// Mark/uncheck sensors as favorites
	MarkSensorFavorite(ctx context.Context, userUuid uuid.UUID, sensorUuid uuid.UUID, favorite bool) error
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
	if sensor.Name == "" || (sensor.Category != domain.SENSOR_CATEGORY_TEMPERATURE && sensor.Category != domain.SENSOR_CATEGORY_PRESSURE && sensor.Category != domain.SENSOR_CATEGORY_HUMIDITY) {
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
	sensor.SensorOwnerUuid = userUuid

	validColors := map[string]bool{
		domain.SENSOR_COLOR_RED:    true,
		domain.SENSOR_COLOR_GREEN:  true,
		domain.SENSOR_COLOR_BLUE:   true,
		domain.SENSOR_COLOR_YELLOW: true,
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

func (s *SensorServiceImpl) EditSensor(ctx context.Context, sensor *domain.Sensor, userUuid uuid.UUID) error {

	var stateOwner, err = s.Repo.GetSensorOwner(ctx, sensor.ID, userUuid)
	if err != nil && !stateOwner {
		return err
	}

	if err = validateRequiredFields(sensor); err != nil {
		return err
	}

	validColors := map[string]bool{
		domain.SENSOR_COLOR_RED:    true,
		domain.SENSOR_COLOR_GREEN:  true,
		domain.SENSOR_COLOR_BLUE:   true,
		domain.SENSOR_COLOR_YELLOW: true,
	}

	if !validColors[sensor.Color] {
		return errors.New("cannot change color if not for one of this: must be RED, GREEN, BLUE, or YELLOW")
	}

	err = s.Repo.EditSensor(ctx, sensor)
	if err != nil {
		return err
	}

	return nil
}

func (s *SensorServiceImpl) ListSensors(ctx context.Context, userUuid uuid.UUID, search string) ([]domain.Sensor, error) {

	var sensors, err = s.Repo.ListSensors(ctx, userUuid, search)
	if err != nil {
		return nil, errors.New("failed to retrieve sensors")
	}

	if search != "" && len(sensors) == 0 {
		return nil, errors.New("no result was found")
	}

	return sensors, nil
}

func (s *SensorServiceImpl) MarkSensorFavorite(ctx context.Context, userUuid uuid.UUID, sensorUuid uuid.UUID, favorite bool) error {
	var err = s.Repo.MarkSensorFavorite(ctx, userUuid, sensorUuid, favorite)
	if err != nil {
		fmt.Print(err)
		return fmt.Errorf("failed to update sensor favorite status")
	}
	return err
}
