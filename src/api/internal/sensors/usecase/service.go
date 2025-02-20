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
	// List all sensors
	ListSensors(ctx context.Context, userUuid uuid.UUID, search string) ([]domain.Sensor, error)
}

// Handles sensor's logic and interaction with the repository
type SensorServiceImpl struct {
	Repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

func (s *SensorServiceImpl) ListSensors(ctx context.Context, userUuid uuid.UUID, search string) ([]domain.Sensor, error) {

	var sensors, err = s.Repo.ListSensors(ctx, userUuid, search)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	if search != "" && len(sensors) == 0 {
		return nil, errors.New("no result was found")
	}

	return sensors, nil
}
