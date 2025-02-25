package usecase

import (
	"api/internal/sensors_data/domain"
	"api/internal/sensors_data/repository"
	"context"
	"fmt"
)

// Interface for sensor's data services
type SensorDataService interface {
	// Add sensor data
	AddSensorData(ctx context.Context, sensorData []*domain.SensorData) error
}

// Handles sensor's data logic and interaction with the repository
type SensorDataServiceImpl struct {
	Repo repository.SensorDataRepository
}

func NewSensorDataService(repo repository.SensorDataRepository) SensorDataService {
	return &SensorDataServiceImpl{Repo: repo}
}

func (s *SensorDataServiceImpl) AddSensorData(ctx context.Context, sensorData []*domain.SensorData) error {
	var err = s.Repo.AddSensorData(ctx, sensorData)
	if err != nil {
		return fmt.Errorf("failed to add sensor data")
	}
	return nil
}
