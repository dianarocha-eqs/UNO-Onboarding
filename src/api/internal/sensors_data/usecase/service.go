package usecase

import (
	"api/internal/sensors_data/domain"
	"api/internal/sensors_data/repository"
	"context"
	"fmt"
	"time"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's data services
type SensorDataService interface {
	// Retrieves sensor data within a specific time interval.
	GetSensorData(ctx context.Context, sensorUuid uuid.UUID, from, to time.Time) ([]domain.SensorData, error)
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

func (s *SensorDataServiceImpl) GetSensorData(ctx context.Context, sensorUuid uuid.UUID, from, to time.Time) ([]domain.SensorData, error) {
	var sensorData, err = s.Repo.GetSensorData(ctx, sensorUuid, from, to)
	if err != nil {
		return sensorData, fmt.Errorf("failed to read sensor data")
	}
	return sensorData, nil
}
