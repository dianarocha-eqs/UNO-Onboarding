package usecase

import (
	"api/internal/sensors/repository"
	"context"
	"fmt"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's services
type SensorService interface {
	MarkSensorFavorite(ctx context.Context, sensorUuid uuid.UUID, favorite bool) error
}

// Handles sensor's logic and interaction with the repository
type SensorServiceImpl struct {
	Repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

func (s *SensorServiceImpl) MarkSensorFavorite(ctx context.Context, sensorUuid uuid.UUID, favorite bool) error {
	var err = s.Repo.MarkSensorFavorite(ctx, sensorUuid, favorite)
	if err != nil {
		return fmt.Errorf("failed to update sensor favorite status")
	}
	return err
}
